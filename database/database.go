//go:generate dbx.v1 golang -d postgres -d sqlite3 -p database memoryLane.dbx .
//go:generate dbx.v1 schema -d postgres -d sqlite3 memoryLane.dbx .

package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

func init() {
	WrapErr = func(e *Error) error {
		return errs.Wrap(e)
	}
	Logger = func(format string, args ...interface{}) {
		fmt.Printf(format+"\n", args...)
	}
}

func (db *DB) WithTx(ctx context.Context,
	fn func(context.Context, *Tx) error) (err error) {

	tx, err := db.Open(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			logrus.Error(err)
			tx.Rollback()
		}
	}()
	return fn(ctx, tx)
}

func joinOR(str ...string) string {
	out := ""
	for i, s := range str {
		if i < len(str)-1 {
			out += "(" + strings.TrimSpace(s) + ")" + " OR "
		} else {
			out += "(" + strings.TrimSpace(s) + ")"
		}
	}

	return out
}

func AppendOrderByAsc(s, col_name string) string {
	return strings.TrimSpace(s) + " ORDER BY " + strings.TrimSpace(col_name) + " ASC;"
}

//lastWeekCitationQuery will return all citations with a memorization date within the last 7 days
//including today
func lastWeekCitationQueryFragment() string {
	return `mem_date BETWEEN (DATE_TRUNC('day', now() - '6 days'::interval)) AND (DATE_TRUNC('day', (now() + '1 day'::interval)))`
}

//dayOfTheWeekQuery returns the db query string that will retreive citations within the last 4
//months memorized on the same day of the week for the last 4 weeks.
func dayOfTheWeekCitationQueryFragment() string {
	return `EXTRACT(ISODOW FROM now()) = EXTRACT(ISODOW FROM mem_date) AND mem_date > (now() - '4.5 weeks'::interval) AND EXTRACT('day' FROM mem_date) != EXTRACT ('day' FROM now())`
}

//dayOfTheMonthQuery will return all citations that have a memorization date on the same day of the
//month as today within the last 12 months.
func dayOfTheMonthCitationQueryFragment() string {
	return `EXTRACT('day' FROM mem_date) = EXTRACT('day' FROM now()) AND  DATE_TRUNC('day', mem_date) != DATE_TRUNC('day', now()) AND mem_date > (now() - '12.5 months'::interval)`
}

func citationSelectorQueryFragment() string {
	return `SELECT reference, author, text, book, hint, year, mem_date FROM citations WHERE `
}

func orderByMemDateFragmentAsc(s string) string {
	return AppendOrderByAsc(s, "mem_date")
}

func allDailyCitationQueryFragments() string {
	return joinOR(
		lastWeekCitationQueryFragment(),
		dayOfTheWeekCitationQueryFragment(),
		dayOfTheMonthCitationQueryFragment(),
	)
}

func buildCompleteDailyCitationQuery() string {
	return orderByMemDateFragmentAsc(
		(citationSelectorQueryFragment() + allDailyCitationQueryFragments()),
	)
}

func (db *DB) GetDailyScriptures() (rows *sql.Rows, err error) {
	return db.Query(buildCompleteDailyCitationQuery())
}
