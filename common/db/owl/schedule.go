package owl

import (
	"fmt"
	"time"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type Schedule struct {
	Name    string
	Timeout int
	Uuid    uuid.UUID
}

func NewSchedule(name string, timeout int) *Schedule {
	return &Schedule{
		Name:    name,
		Timeout: timeout,
	}
}

type OwlSchedule struct {
	Id             int       `db:"sch_id"`
	Name           string    `db:"sch_name"`
	Lock           byte      `db:"sch_lock"`
	LastUpdateTime time.Time `db:"sch_modify_time"`
}

func (sch *OwlSchedule) isLocked() bool {
	return sch.Lock == byte(LOCKED)
}

type OwlScheduleLog struct {
	Uuid      cdb.DbUuid `db:"sl_uuid"`
	SchId     int        `db:"sl_sch_id"`
	StartTime time.Time  `db:"sl_start_time"`
	EndTime   time.Time  `db:"sl_end_time"`
	Timeout   int        `db:"sl_timeout"`
	Status    byte       `db:"sl_status"`
	Message   string     `db:"sl_message"`
}

var insertSql = `
	INSERT INTO owl_schedule_log(
		sl_uuid, sl_sch_id,
		sl_start_time, sl_timeout, sl_status
	)
	VALUES (:uuid, :schid, :starttime, :timeout, :status)
`

type TimeOutOfSchedule struct {
	Name          string
	AcquiredTime  time.Time
	LastStartTime time.Time
	Timeout       int
}

func (t *TimeOutOfSchedule) Error() string {
	return fmt.Sprintf("%s error: period between %s and %s should longer than %ds", t.Name,
		t.LastStartTime.Format(time.RFC3339),
		t.AcquiredTime.Format(time.RFC3339),
		t.Timeout)
}

type LockStatus byte

const (
	FREE LockStatus = iota
	LOCKED
)

type TaskStatus byte

const (
	DONE TaskStatus = iota
	RUN
	FAIL
	TIMEOUT
)

func AcquireLock(schedule *Schedule) error {
	txProcessor := &txAcquireLock{
		schedule:     schedule,
		timeoutError: nil,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
	return txProcessor.timeoutError
}

type txAcquireLock struct {
	schedule     *Schedule
	timeoutError *TimeOutOfSchedule

	lockTable OwlSchedule
	logTable  OwlScheduleLog
}

func (ack *txAcquireLock) InTx(tx *sqlx.Tx) cdb.TxFinale {

	now := time.Now()

	/**
	 * Lock table
	 */
	ack.selectOrInsertLock(tx, now)
	// The previous task is not timeout()
	if ack.lockTable.isLocked() && ack.notTimeout(tx, now) {
		ack.timeoutError = &TimeOutOfSchedule{
			Name:          "Schedule locked",
			LastStartTime: ack.logTable.StartTime,
			AcquiredTime:  now,
			Timeout:       ack.logTable.Timeout,
		}
		return cdb.TxCommit
	}

	ack.updateLockByName(tx, now)
	// :~)

	/**
	 * Log table
	 */
	generatedUuid := uuid.NewV4()
	_ = sqlxExt.ToTxExt(tx).NamedExec(insertSql,
		map[string]interface{}{
			"uuid":      cdb.DbUuid(generatedUuid),
			"schid":     ack.lockTable.Id,
			"starttime": now,
			"timeout":   ack.schedule.Timeout,
			"status":    RUN,
		},
	)
	ack.schedule.Uuid = generatedUuid
	// :~)

	return cdb.TxCommit
}

func (ack *txAcquireLock) selectOrInsertLock(tx *sqlx.Tx, now time.Time) {
	name := ack.schedule.Name
	exist := sqlxExt.ToTxExt(tx).GetOrNoRow(&ack.lockTable, `
		SELECT sch_id, sch_lock
		FROM owl_schedule
		WHERE sch_name = ?
		FOR UPDATE
	`, name)

	if !exist {
		r := tx.MustExec(`
			INSERT INTO owl_schedule(
				sch_name,
				sch_lock, sch_modify_time
			)
			VALUES (?, 0, ?)
		`, name, now)
		ack.lockTable.Id = int(cdb.ToResultExt(r).LastInsertId())
		ack.lockTable.Lock = byte(FREE)
	}
}

func (ack *txAcquireLock) updateLockByName(tx *sqlx.Tx, now time.Time) {
	_ = tx.MustExec(`
		UPDATE owl_schedule
		SET sch_lock = 1,
			sch_modify_time = ?
		WHERE sch_name = ?
	`, now, ack.schedule.Name)
}

func (ack *txAcquireLock) notTimeout(tx *sqlx.Tx, now time.Time) bool {
	ret := &ack.logTable
	exist := sqlxExt.ToTxExt(tx).GetOrNoRow(ret, `
		SELECT sl.sl_start_time, sl.sl_timeout
		FROM owl_schedule_log sl
		WHERE sl.sl_sch_id = ?
		ORDER BY sl.sl_start_time DESC
		LIMIT 1
	`, ack.lockTable.Id)

	// Check timeout iff row exists
	return exist && (now.Sub(ret.StartTime) <= time.Duration(ret.Timeout)*time.Second)
}
