package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Eun/logtimer"
	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	var relativeFlag string
	var formatFlag string
	var colorCorrection string
	var rootCmd = &cobra.Command{
		Use: filepath.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			var cc logtimer.ColorCorrection
			switch strings.ToLower(colorCorrection) {
			case "true", "normal", "standard", "enable":
				cc = logtimer.Enabled
			case "alternate":
				cc = logtimer.Alternate
			default:
				cc = logtimer.Disabled
			}

			var format logtimer.FormatFunc

			if relativeFlag != "" {
				startTime := time.Now()
				format = func() string {
					return logtimer.FormatDuration(time.Since(startTime), relativeFlag)
				}
			} else {
				format = func() string {
					return logtimer.FormatTime(time.Now(), formatFlag)
				}
			}

			reader := &logtimer.PrefixReader{
				Reader:          os.Stdin,
				Format:          format,
				ColorCorrection: cc,
			}

			_, _ = io.Copy(os.Stdout, reader)
		},
	}
	rootCmd.Version = version + " " + date + " " + commit
	rootCmd.Flags().StringVarP(&formatFlag, "format", "f", "[%X] ", `format to prefix the lines. You can use following directives to format the date:
	%a    Weekday as locale’s abbreviated name.                             (Sun, Mon, ..., Sat)
	%A    Weekday as locale’s full name.                                    (Sunday, Monday, ..., Saturday)
	%w    Weekday as a decimal number, where 0 is Sunday and 6 is Saturday  (0, 1, ..., 6)
	%d    Day of the month as a decimal number.                             (1, 2, ..., 31)
	%b    Month as locale’s abbreviated name.                               (Jan, Feb, ..., Dec)
	%B    Month as locale’s full name.                                      (January, February, ..., December)
	%m    Month as a decimal number.                                        (1, 2, ..., 12)
	%y    Year without century as a decimal number.                         (0, 1, ..., 99)
	%Y    Year with century as a decimal number.                            (1970, 1988, 2001, 2013)
	%H    Hour (24-hour clock) as a decimal number.                          (0, 1, ..., 23)
	%I    Hour (12-hour clock) as a decimal number.                          (1, 2, ..., 12)
	%p    Meridian indicator.                                               (AM, PM)
	%M    Minute as a decimal number.                                       (0, 1, ..., 59)
	%S    Second as a decimal number.                                       (0, 1, ..., 59)
	%f    Microsecond as a decimal number.                                  (0, 1, ..., 999999)
	%z    UTC offset in the form +HHMM or -HHMM                             (+0000)
	%Z    Time zone name                                                    (UTC)
	%j    Day of the year as a decimal number                               (1, 2, ..., 366)
	%U    Week number of the year (Sunday as the first day of the week) as a decimal number. All days in a new year preceding the first Sunday are considered to be in week 0.
	                                                                        (0, 1, ..., 53)
	%W    Week number of the year (Monday as the first day of the week) as a decimal number. All days in a new year preceding the first Monday are considered to be in week 0.
	                                                                        (0, 1, ..., 53)
	%c    Date and time representation.                                     (Tue Aug 16 21:30:00 1988)
	%x    Date representation.                                              (08/16/88)
	%X    Time representation.                                              (21:30:00)
	%%    A literal '%' character.                                          (%)

	It is possible to to zero/space pad the directives
    Example:
		ping 8.8.8.8 | logtimer --format="[%a, %d %b %Y %02H:%02M:%02S %Z] "

`)
	rootCmd.Flags().StringVarP(&relativeFlag, "relative", "r", "", `use relative log mode, this means that the clock will start at execution date. You can use following directives to format the time
	%X    Total Time elapsed.                                               (85:30:04)
	%Xf   Total Time with Microseconds elapsed.                             (85:30:04.999999)
	%Xn   Total Time with Nanoseconds elapsed.                              (85:30:04.999999999)
	%%    A literal '%' character.                                          (%)

	It is possible to to zero/space pad the directives
	Examples:
		$ ping 8.8.8.8 | logtimer --relative="[%010X] "
		[0000:00:00] 64 bytes from 8.8.8.8: icmp_seq=21 ttl=123 time=18.3 ms
		...
		[0085:30:04] 64 bytes from 8.8.8.8: icmp_seq=22 ttl=123 time=18.5 ms
		...
		[9854:30:04] 64 bytes from 8.8.8.8: icmp_seq=22 ttl=123 time=18.5 ms
	`)
	rootCmd.Flag("relative").NoOptDefVal = "[%X] "

	rootCmd.Flags().StringVarP(&colorCorrection, "color-correction", "", "enable", "change color correction if you experience problems (possible values: enable, alternate, disable")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
