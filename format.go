package logtimer

import (
	"fmt"
	"math"
	"time"

	"github.com/Eun/mapprint"
)

var longDayNames = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
}

var shortMonthNames = []string{
	"---",
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var longMonthNames = []string{
	"---",
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func weekNumber(t time.Time, char int) int {
	weekday := int(t.Weekday())

	if char == 'W' {
		// Monday as the first day of the week
		if weekday == 0 {
			weekday = 6
		} else {
			weekday--
		}
	}

	return (t.YearDay() + 6 - weekday) / 7
}

// FormatTime formats a time, use the follwing format:
// %a    Weekday as locale’s abbreviated name.                             (Sun, Mon, ..., Sat)
// %A    Weekday as locale’s full name.                                    (Sunday, Monday, ..., Saturday)
// %w    Weekday as a decimal number, where 0 is Sunday and 6 is Saturday  (0, 1, ..., 6)
// %d    Day of the month as a decimal number.                             (1, 2, ..., 31)
// %b    Month as locale’s abbreviated name.                               (Jan, Feb, ..., Dec)
// %B    Month as locale’s full name.                                      (January, February, ..., December)
// %m    Month as a decimal number.                                        (1, 2, ..., 12)
// %y    Year without century as a decimal number.                         (0, 1, ..., 99)
// %Y    Year with century as a decimal number.                            (1970, 1988, 2001, 2013)
// %H    Hour (24-hour clock) as a decimal number.                          (0, 1, ..., 23)
// %I    Hour (12-hour clock) as a decimal number.                          (1, 2, ..., 12)
// %p    Meridian indicator.                                               (AM, PM)
// %M    Minute as a decimal number.                                       (0, 1, ..., 59)
// %S    Second as a decimal number.                                       (0, 1, ..., 59)
// %f    Microsecond as a decimal number.                                  (0, 1, ..., 999999)
// %z    UTC offset in the form +HHMM or -HHMM                             (+0000)
// %Z    Time zone name                                                    (UTC)
// %j    Day of the year as a decimal number                               (1, 2, ..., 366)
// %U    Week number of the year (Sunday as the first day of the week) as a decimal number. All days in a new year preceding the first Sunday are considered to be in week 0.
//
//	(0, 1, ..., 53)
//
// %W    Week number of the year (Monday as the first day of the week) as a decimal number. All days in a new year preceding the first Monday are considered to be in week 0.
//
//	(0, 1, ..., 53)
//
// %c    Date and time representation.                                     (Tue Aug 16 21:30:00 1988)
// %x    Date representation.                                              (08/16/88)
// %X    Time representation.                                              (21:30:00)
// %%    A literal '%' character.                                          (%)
func FormatTime(t time.Time, f string) string {
	return mapprint.Sprintf(f, map[string]interface{}{
		"a": shortDayNames[t.Weekday()],
		"A": longDayNames[t.Weekday()],
		"w": t.Weekday,
		"d": t.Day,
		"b": shortMonthNames[t.Month()],
		"B": longMonthNames[t.Month()],
		"m": t.Month,
		"y": t.Year() % 100,
		"Y": t.Year,
		"H": t.Hour,
		"I": func() int {
			if t.Hour() == 0 {
				return 12
			} else if t.Hour() > 12 {
				return t.Hour() - 12
			} else {
				return t.Hour()
			}
		},
		"p": func() string {
			if t.Hour() < 12 {
				return "AM"
			}
			return "PM"
		},
		"M": t.Minute,
		"S": t.Second,
		"f": t.Nanosecond() / 1000,
		"z": t.Format("-0700"),
		"Z": t.Format("MST"),
		"j": t.YearDay(),
		"U": weekNumber(t, 'U'),
		"W": weekNumber(t, 'W'),
		"c": t.Format("Mon Jan 2 15:04:05 2006"),
		"x": fmt.Sprintf("%02d/%02d/%02d", t.Month(), t.Day(), t.Year()%100),
		"X": fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()),
	})
}

// FormatDuration formats a duration, use the follwing format:
// %X    Total Time elapsed.                                              (85:30:04)
// %Xf   Total Time with Microseconds elapsed.                            (85:30:04.999999)
// %Xn   Total Time with Nanoseconds elapsed.                             (85:30:04.999999999)
// %%    A literal '%' character.                                         (%)
func FormatDuration(d time.Duration, f string) string {
	return mapprint.Sprintf(f, map[string]interface{}{
		"X": func() string {
			s := d.Seconds()
			hours := int(math.Floor(s / 3600))
			//nolint: staticcheck // calling math.Floor on a converted integer is pointless (staticcheck)
			minutes := int(math.Floor(float64(int(s) % 3600 / 60)))
			//nolint: staticcheck // calling math.Floor on a converted integer is pointless (staticcheck)
			seconds := int(math.Floor(float64(int(s) % 3600 % 60)))
			return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		},
		"Xf": func() string {
			return FormatDuration(d, "%X") + fmt.Sprintf(".%06d", d%time.Second/time.Microsecond)
		},
		"Xn": func() string {
			return FormatDuration(d, "%X") + fmt.Sprintf(".%09d", d%time.Second)
		},
	})
}
