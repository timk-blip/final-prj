package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}
	dstartdate, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", err
	}
	part := strings.Split(repeat, " ")
	switch part[0] {
	case "y":
		for {
			dstartdate = dstartdate.AddDate(1, 0, 0)
			if afterNow(dstartdate, now) {
				break
			}
		}
		return dstartdate.Format(DateLayout), nil
	case "d":
		if len(part) < 2 {
			return "", errors.New("day is empty")
		}
		interval, err := strconv.Atoi(part[1])
		if err != nil {
			return "", err
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("interval must be between 1 and 400")
		}
		for {
			dstartdate = dstartdate.AddDate(0, 0, interval)
			if afterNow(dstartdate, now) {
				break
			}
		}
		return dstartdate.Format(DateLayout), nil
	case "w":
		var week [8]bool
		if len(part) != 2 {
			return "", errors.New("invalid format day of the week")
		}
		partWeek := strings.Split(part[1], ",")
		for _, s := range partWeek {
			n, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if n < 1 || n > 7 {
				return "", errors.New("the day of the week must be between 1 and 7")
			}
			week[n] = true
		}
		for {
			dstartdate = dstartdate.AddDate(0, 0, 1)
			day := int(dstartdate.Weekday())
			if day == 0 {
				day = 7
			}
			if afterNow(dstartdate, now) && week[day] {
				break
			}
		}
		return dstartdate.Format(DateLayout), nil
	case "m":
		if len(part) != 2 && len(part) != 3 {
			return "", errors.New("invalid format for m")
		}
		var monthS [13]bool
		var dayS [32]bool
		lastDay, penultDay := false, false
		partDays := strings.Split(part[1], ",")
		for _, s := range partDays {
			n, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if n == -1 {
				lastDay = true
			} else if n == -2 {
				penultDay = true
			} else if n >= 1 && n <= 31 {
				dayS[n] = true
			} else {
				return "", errors.New("the day of the month must be between 1 and 31")
			}
		}
		monthsGiven := len(part) == 3
		if monthsGiven {
			partMonth := strings.Split(part[2], ",")
			for _, s := range partMonth {
				n, err := strconv.Atoi(s)
				if err != nil {
					return "", err
				}
				if n >= 1 && n <= 12 {
					monthS[n] = true
				} else {
					return "", errors.New("the month of the year be between 1 and 12")
				}
			}
		}
		for {
			dstartdate = dstartdate.AddDate(0, 0, 1)
			day := dstartdate.Day()
			daysInMonth := time.Date(dstartdate.Year(), dstartdate.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Day()
			dayMatch := dayS[day] || (lastDay && day == daysInMonth) || (penultDay && day == daysInMonth-1)
			month := int(dstartdate.Month())
			monthMatch := !monthsGiven || monthS[month]
			if afterNow(dstartdate, now) && dayMatch && monthMatch {
				break
			}
		}
		return dstartdate.Format(DateLayout), nil

	default:
		return "", errors.New("invalid repeat")
	}
}

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

func nextdateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now := time.Now()
	if nowParam != "" {
		var err error
		now, err = time.Parse(DateLayout, nowParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}
	}
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, nextDate)
}
