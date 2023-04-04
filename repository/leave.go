package repository

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type Leave struct{}

func (leave *Leave) GetEmployInformation(input model.Leave) dto.ResponseDto {
	var op dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output dto.LeaveEmpInfo

	result := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate ,c.deptname ,d.designame 
	from payroll.staffinformation a
	inner join payroll.salarystructure b 
	on a.empcode =b.empcode 
	INNER JOIN payroll.hrdept c
	ON c.deptcode =b.deptcode
	INNER JOIN payroll.designation d
	ON d.desigcode =b.desigcode where a.empcode=?`, input.Empcode).First(&output)
	output.Confirmdate = output.Confirmdate[:10]
	output.Joindate = output.Joindate[:10]
	if result.RowsAffected == 0 {
		op.IsSuccess = false
		op.Message = "Empcode not found"
		op.StatusCode = http.StatusNotFound
		return op
	}
	op.IsSuccess = true
	op.Message = "Success"
	op.StatusCode = 200
	op.Payload = output
	return op
}

func (leave *Leave) CheckConsumeLeave(input dto.LeaveCount) dto.ResponseDto {
	var op dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output dto.LeaveCount
	var output3 []dto.LeaveCount
	var staff model.Staffinformation
	var output2 dto.LeaveEmpInfo
	var holidayop model.Holidayinfo
	var weekendop model.Weekend
	var leaverecordinput model.Leaverecord
	var leaverecordoutput model.Leaverecord
	var leavecheck model.Leaverecord

	date := time.Now()
	currentdate := date.Format("2006-01-02")
	currentyear := date.Year()
	fmt.Println("Current time", date)
	fmt.Println("current year", currentyear)
	startdate := strconv.Itoa(currentyear) + "-01-01"
	enddate := strconv.Itoa(currentyear) + "-12-31"

	// result3 := db.Raw("select * from payroll.leaverecord where year=?", currentyear).First(&leavecheck)
	// if result3.RowsAffected == 0 {
	// 	op.IsSuccess = false
	// 	op.StatusCode = 404
	// 	op.Message = "Leave record not found. Please update carry forward leaves"
	// 	op.Payload = nil
	// 	return op
	// } else {
		// casual leave
		if input.Lcode == 1 {

			result := db.Raw("select * from payroll.leavetype where lcode = ?", input.Lcode).First(&output)
			// result = db.Where("lcode = ?", input.Lcode).First(&model.Leavetype{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Type not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select casual from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofcasual := 0
			for i := 0; i < len(output3); i++ {
				sumofcasual = sumofcasual + output3[i].Casual
			}
			output.Casual = sumofcasual
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
						from payroll.staffinformation a
						inner join payroll.salarystructure b
						on a.empcode =b.empcode
						where a.empcode=?`, input.Empcode).First(&output2)
			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Emp not found"
				op.Payload = nil
				return op
			}
			if output2.Confirmdate == "" {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Confirm date not found"
				op.Payload = nil
				return op
			}

			// medical leave
		} else if input.Lcode == 2 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}

			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result := db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select medical from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }

			sumofmedicalleave := 0
			for i := 0; i < len(output3); i++ {
				sumofmedicalleave = sumofmedicalleave + output3[i].Medical
			}

			output.Medical = sumofmedicalleave

			// earn leave
		} else if input.Lcode == 3 {
			result4 := db.Raw("select * from payroll.leaverecord where year=?", currentyear).First(&leavecheck)
			if result4.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave record not found. Please update carry forward leaves"
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select earn from payroll.leave where empcode =?", input.Empcode).First(&output)

			// _ = db.Raw("select * from payroll.leave where empcode=? and leavedate=?", input.Empcode, input.Leavedate).Find(&output3)

			// sumofearnleave := 0
			// sumofallleave := 0
			// for i := 0; i < len(output3); i++ {

			// 	sumofearnleave = sumofearnleave + output3[i].Earn
			// 	sumofallleave = (sumofallleave + output3[i].Casual + output3[i].Medical + output3[i].Special + output3[i].Earn + output3[i].Study +
			// 		output3[i].Maternity + output3[i].Leavewopay + output3[i].Extraordi + output3[i].Festival + output3[i].Sick)

			// }
			// output.Earn = sumofearnleave

			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			if output2.Confirmdate == "" {
				op.IsSuccess = false
				op.Message = "Exception! Employee is not confirmed yet"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			} else if output2.Confirmdate[:4] == currentdate[:4] {
				fmt.Println("Current Date", currentdate)
				// fmt.Println("Joining Date", output2.Joindate[:10])
				fmt.Println("confirm Date", output2.Confirmdate)

				cdate := currentdate[:4]
				lastdate := cdate + "-12-31"

				fmt.Println("Last Date- ", lastdate)
				format := "2006-01-02"
				then, _ := time.Parse(format, lastdate)

				fmt.Println("after time parse current date : ", then)

				confirmdate := output2.Confirmdate[:10]
				fmt.Println("Before parsing confirm date ", confirmdate)
				now, _ := time.Parse(format, confirmdate)

				fmt.Println("after time parse confirmation date: ", now)

				diff := then.Sub(now)
				days := (int(diff.Hours() / 24)) + 1
				fmt.Println("days : ", days)

				_ = db.Raw("select count(noofdays) noofdays from payroll.holidayinfo where year=?", cdate).First(&holidayop)
				_ = db.Raw("select count(noofweek) noofweek from payroll.weekend where year=?", cdate).First(&weekendop)
				firstdate := currentdate[:4]
				firstdate = firstdate + "-01-01"
				_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ? ", input.Empcode, firstdate, currentdate).Find(&output3)

				sumofearnleave := 0
				sumofallleave := 0
				for i := 0; i < len(output3); i++ {

					sumofearnleave = 0
					sumofallleave = (sumofallleave + output3[i].Casual + output3[i].Medical + output3[i].Special + output3[i].Earn + output3[i].Study +
						output3[i].Maternity + output3[i].Leavewopay + output3[i].Extraordi + output3[i].Festival + output3[i].Sick)

				}
				output.Earn = sumofearnleave

				fmt.Println("sum of all leave ", sumofallleave)
				fmt.Println("sum of earn leave", sumofearnleave)
				fmt.Println("sum of holiday leave ", holidayop.Noofdays)
				fmt.Println("sum of weekend leave ", weekendop.Noofweek)
				minus := sumofallleave + holidayop.Noofdays + weekendop.Noofweek
				fmt.Println("minus days", minus)
				earnable_days := days - minus
				fmt.Println("earnable_days - ", earnable_days)
				actual_earnleave := earnable_days / 11
				fmt.Println("actual_earnleave - ", actual_earnleave)

				leaveyrs := output2.Confirmdate[:4]
				fmt.Println("leaveyrs - ", leaveyrs)
				// leaverecordinput.El_earned_nextyear = actual_earnleave
				leaverecordinput.Empcode = input.Empcode
				z, _ := strconv.Atoi(leaveyrs)
				leaverecordinput.Year = z

				// if already calculated then update the input of leave record
				result1 := db.Raw("select * from payroll.leaverecord where empcode = ? and year = ? ", input.Empcode, leaverecordinput.Year).First(&leaverecordoutput)
				if result1.RowsAffected >= 1 {
					leaverecordinput.Year = leaverecordoutput.Year
					leaverecordinput.Empcode = leaverecordoutput.Empcode
					leaverecordoutput.El_earned_nextyear = actual_earnleave
					leaverecordoutput.Earnconsume = 0
					leaverecordoutput.Earnleave = 0
					leaverecordoutput.Earnconsume = 0
					leaverecordoutput.Rest_of_leave = 0
					_ = db.Where(&model.Leaverecord{Id: leaverecordinput.Id, Empcode: leaverecordinput.Empcode, Year: leaverecordinput.Year}).Updates(&leaverecordoutput)

				} else {
					_ = db.Raw("select coalesce ((max(id)+1),1) from payroll.leaverecord").First(&leaverecordinput.Id)
					leaverecordinput.El_earned_nextyear = actual_earnleave
					leaverecordinput.Earnconsume = 0
					leaverecordinput.Earnleave = 0
					leaverecordinput.Rest_of_leave = 0

					_ = db.Create(&leaverecordinput)
					// if result.RowsAffected == 0 {
					// 	op.IsSuccess = false
					// 	op.Message = "failed to create leave record"
					// 	op.Payload = leaverecordinput
					// 	op.StatusCode = 404
					// 	return op
					// }
				}
				// if not calculated previously then create a new one of leave record

				output.Earn = actual_earnleave
				op.IsSuccess = true
				op.Message = "You can't enojoy earn leave this year"
				op.StatusCode = 200
				op.Payload = output
				return op

			} else if output2.Confirmdate[:4] != currentdate[:4] {

				fmt.Println("Current Date", currentdate)

				c_year := currentdate[:4]
				firstday := c_year + "-01-01"
				lastday := c_year + "-12-31"

				fmt.Println("Last Date ", lastday)
				format := "2006-01-02"
				enddate, _ := time.Parse(format, lastday)

				fmt.Println("after time parse last date : ", enddate)

				// b := output2.Confirmdate[:10]
				fmt.Println("Before parsing first date ", firstday)
				startdate, _ := time.Parse(format, firstday)

				fmt.Println("after time parse first date: ", startdate)

				diff := enddate.Sub(startdate)
				days := (int(diff.Hours() / 24)) + 1
				fmt.Println("days : ", days)

				_ = db.Raw("select count(noofdays) noofdays from payroll.holidayinfo where year=?", c_year).First(&holidayop)
				_ = db.Raw("select count(noofweek) noofweek from payroll.weekend where year=?", c_year).First(&weekendop)
				_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ? ", input.Empcode, firstday, lastday).Find(&output3)

				sumofearnleave := 0
				sumofallleave := 0
				for i := 0; i < len(output3); i++ {

					sumofearnleave = sumofearnleave + output3[i].Earn
					sumofallleave = (sumofallleave + output3[i].Casual + output3[i].Medical + output3[i].Special + output3[i].Earn + output3[i].Study +
						output3[i].Maternity + output3[i].Leavewopay + output3[i].Extraordi + output3[i].Festival + output3[i].Sick)

				}
				output.Earn = sumofearnleave

				fmt.Println("sum of all leave ", sumofallleave)
				fmt.Println("sum of earn leave", sumofearnleave)
				fmt.Println("sum of holiday leave ", holidayop.Noofdays)
				fmt.Println("sum of weekend leave ", weekendop.Noofweek)
				minus := sumofallleave + holidayop.Noofdays + weekendop.Noofweek
				fmt.Println("minus days", minus)
				earnable_days := days - minus
				fmt.Println("earnable_days - ", earnable_days)
				actual_earnleave := earnable_days / 11
				fmt.Println("actual_earnleave - ", actual_earnleave)

				leaveyrs := currentdate[:4]
				fmt.Println("leaveyears in string- ", leaveyrs)
				leaverecordinput.El_earned_nextyear = actual_earnleave
				leaverecordinput.Empcode = input.Empcode
				z, _ := strconv.Atoi(leaveyrs)
				leaverecordinput.Year = z
				fmt.Println("Leave year ", leaverecordinput.Year)

				// if already calculated then update the input of leave record
				result1 := db.Raw("select * from payroll.leaverecord where empcode = ? and year = ? ", input.Empcode, leaverecordinput.Year).First(&leaverecordoutput)
				if result1.RowsAffected >= 1 {
					output.Earndays = leaverecordoutput.Earnleave
					output.Earn = sumofearnleave
					output.Empcode = leaverecordoutput.Empcode

					// _ = db.Raw("select * from payroll.leaverecord where empcode = ? and year = ? ", input.Empcode, leaverecordinput.Year-1).First(&leaverecordoutput)
					// leaverecordoutput.Year = leaverecordinput.Year
					// leaverecordoutput.Empcode = leaverecordinput.Empcode

					//leaverecordinput.Earnleave = leaverecordoutput.El_earned_nextyear + leaverecordoutput.Rest_of_leave
					if leaverecordoutput.Earnleave == 0 {
						leaverecordinput.Earnconsume = 0
					} else {
						leaverecordinput.Earnconsume = sumofearnleave
					}

					leaverecordinput.El_earned_nextyear = actual_earnleave
					//leaverecordinput.Earnconsume = sumofearnleave
					leaverecordinput.Rest_of_leave = leaverecordoutput.Earnleave - leaverecordinput.Earnconsume

					_ = db.Where(&model.Leaverecord{Id: leaverecordinput.Id, Empcode: leaverecordinput.Empcode, Year: leaverecordinput.Year}).Updates(&leaverecordinput)
					// output.Earndays = leaverecordinput.Earnleave

				}
				// else {
				// 	_ = db.Raw("select * from payroll.leaverecord where empcode = ? and year = ? ", input.Empcode, leaverecordinput.Year-1).First(&leaverecordoutput)
				// 	_ = db.Raw("select coalesce ((max(id)+1),1) from payroll.leaverecord").First(&leaverecordinput.Id)

				// 	// leaverecordinput.Year = leaverecordinput.Year

				// 	leaverecordinput.Earnleave = leaverecordoutput.El_earned_nextyear + leaverecordoutput.Rest_of_leave
				// 	if leaverecordoutput.Earnleave == 0 {
				// 		leaverecordinput.Earnconsume = 0
				// 	} else {
				// 		leaverecordinput.Earnconsume = sumofearnleave
				// 	}
				// 	leaverecordinput.Rest_of_leave = leaverecordinput.Earnleave - leaverecordinput.Earnconsume
				// 	leaverecordinput.El_earned_nextyear = actual_earnleave

				// 	fmt.Println("Earn Leave from pre year", leaverecordinput.Earnleave)
				// 	// if leaverecordinput.El_earned_nextyear > 60 {
				// 	// 	leaverecordinput.El_earned_nextyear = 60
				// 	// }
				// 	_ = db.Create(&leaverecordinput)
				// 	output.Earndays = leaverecordinput.Earnleave
				// }

			}

			// maternity leave
		} else if input.Lcode == 4 {
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select gender from payroll.staffinformation where empcode=?", input.Empcode).First(&staff)
			// if result3.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Staff Information not found"
			// 	op.Payload = nil
			// 	return op
			// }
			if staff.Gender == "male" {
				op.IsSuccess = false
				op.StatusCode = 400
				op.Message = "Male Gender is not allowed for maternity leave"
				op.Payload = nil
				return op
			}

			_ = db.Raw("select maternity from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofmaternityleave := 0
			for i := 0; i < len(output3); i++ {
				sumofmaternityleave = sumofmaternityleave + output3[i].Maternity
			}
			output.Maternity = sumofmaternityleave

			// festical leave
		} else if input.Lcode == 5 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select festival from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumoffestivalleave := 0
			for i := 0; i < len(output3); i++ {
				sumoffestivalleave = sumoffestivalleave + output3[i].Festival
			}
			output.Festival = sumoffestivalleave

			// sick leave
		} else if input.Lcode == 6 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select sick from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofsickleave := 0
			for i := 0; i < len(output3); i++ {
				sumofsickleave = sumofsickleave + output3[i].Sick
			}
			output.Sick = sumofsickleave
			// extra ordinary leave
		} else if input.Lcode == 7 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select extraordi from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofextraordileave := 0
			for i := 0; i < len(output3); i++ {
				sumofextraordileave = sumofextraordileave + output3[i].Extraordi
			}
			output.Extraordi = sumofextraordileave
			// study leave
		} else if input.Lcode == 8 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select study from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofstudyleave := 0
			for i := 0; i < len(output3); i++ {
				sumofstudyleave = sumofstudyleave + output3[i].Study
			}
			output.Study = sumofstudyleave
			// leave with out pay
		} else if input.Lcode == 9 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select leavewopay from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofleavewopayleave := 0
			for i := 0; i < len(output3); i++ {
				sumofleavewopayleave = sumofleavewopayleave + output3[i].Leavewopay
			}
			output.Leavewopay = sumofleavewopayleave

			// special leave
		} else if input.Lcode == 10 {
			result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
							from payroll.staffinformation a
							inner join payroll.salarystructure b
							on a.empcode =b.empcode
							where a.empcode=?`, input.Empcode).First(&output2)

			if result3.RowsAffected == 0 {
				op.IsSuccess = false
				op.Message = "Exception! Empcode not found"
				op.StatusCode = http.StatusNotFound
				op.Payload = nil
				return op
			}
			result := db.Raw("select * from payroll.leavetype where lcode =?", input.Lcode).First(&output)
			// result = db.Where("lcode =?", input.Lcode).First(&model.LeaveCount{})
			if result.RowsAffected == 0 {
				op.IsSuccess = false
				op.StatusCode = 404
				op.Message = "Leave Count not found"
				op.Payload = nil
				return op
			}
			_ = db.Raw("select special from payroll.leave where empcode =?", input.Empcode).First(&output)
			// if result2.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = 404
			// 	op.Message = "Leave not found"
			// 	op.Payload = nil
			// 	return op
			// }
			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, startdate, enddate).Find(&output3)
			// if result1.RowsAffected == 0 {
			// 	op.IsSuccess = false
			// 	op.StatusCode = http.StatusNotFound
			// 	op.Message = "Leave not found"
			// 	return op
			// }
			sumofspecialleave := 0
			for i := 0; i < len(output3); i++ {
				sumofspecialleave = sumofspecialleave + output3[i].Special
			}
			output.Special = sumofspecialleave
		} else {
			op.IsSuccess = false
			op.StatusCode = 404
			op.Message = "Leave not found"
			op.Payload = nil
			return op
		}
	
	// result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
	// from payroll.staffinformation a
	// inner join payroll.salarystructure b
	// on a.empcode =b.empcode
	// where a.empcode=?`, input.Empcode).First(&output2)

	// if result3.RowsAffected == 0 {
	// 	op.IsSuccess = false
	// 	op.Message = "Empcode not found"
	// 	op.StatusCode = http.StatusNotFound
	// 	op.Payload = nil
	// 	return op
	// }

	// fmt.Println("Joining Date", output2.Joindate)
	// // fmt.Println("Joining Date", output2.Joindate[:10])
	// fmt.Println("confirm Date", output2.Confirmdate)
	// // a := output2.Joindate
	// a := output2.Joindate[:10]
	// date := time.Now()
	// fmt.Println(" Date - ", date)
	// format := "2006-01-02"
	// then, _ := time.Parse(format, a)

	// fmt.Println("after time parse joining date : ", then)

	// // b := output2.Confirmdate
	// b := output2.Confirmdate[:10]
	// now, _ := time.Parse(format, b)

	// fmt.Println("after time parse confirmation date: ", now)

	// diff := date.Sub(then)
	// days := (int(diff.Hours() / 24))
	// fmt.Println("days : ", days)

	op.IsSuccess = true
	op.Message = "Success"
	op.StatusCode = 200
	op.Payload = output
	return op
}

func (leave *Leave) EntryANewLeave(input model.Leave) dto.ResponseDto {
	var op dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint")

	var output []model.Leave
	var leaverecord model.Leaverecord

	if input.Empcode <= 0 {
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Message = "Empcode is required"
		return op
	}
	date := time.Now()
	cdate := date.Format("2006-01-02")
	year := cdate[:4]
	fmt.Println("year", year)

	leave_date := input.Leavedate[:4]
	if leave_date != year {
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Message = "Leave date is not matched with current year!"
		return op
	}
	leave_from := input.Leavefrom[:4]
	if leave_from != year {
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Message = "Leave from date is not matched with current year!"
		return op
	}
	leave_to := input.Leaveto[:4]
	if leave_to != year {
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Message = "Leave to date is not matched with current year!"
		return op
	}

	if input.Casual <= 0 && input.Earn <= 0 && input.Extraordi <= 0 && input.Festival <= 0 && input.Leavewopay <= 0 && input.Maternity <= 0 && input.Medical <= 0 && input.Sick <= 0 && input.Special <= 0 && input.Study <= 0 {
		tx.RollbackTo("savepoint")
		op.Message = "Insert any days leave before entry"
		op.IsSuccess = false
		op.Payload = nil
		return op
	}
	if input.Cause == "" {
		tx.RollbackTo("savepoint")
		op.Message = "Write cause of leave before leave entry"
		op.IsSuccess = false
		op.Payload = nil
		return op

	}
	format := "2006-01-02"
	f_date, _ := time.Parse(format, input.Leavefrom[:10])
	t_date, _ := time.Parse(format, input.Leaveto[:10])
	diff := t_date.Sub(f_date)
	days := (int(diff.Hours() / 24)) + 1
	fmt.Println("How many days : ", days)

	sumofallleave := (input.Casual + input.Earn + input.Extraordi + input.Festival +
		input.Leavewopay + input.Maternity + input.Medical + input.Sick +
		input.Special + input.Study)

	a := "Your leave dates are not permit more than leaves leave you enter"
	if sumofallleave != days {
		tx.RollbackTo("savepoint")
		op.IsSuccess = false
		op.StatusCode = 404
		op.Message = a
		op.Payload = nil
		return op

	}

	result1 := tx.Raw("select * from payroll.leaverecord where empcode=? and year=?", input.Empcode, year).First(&leaverecord)
	if result1.RowsAffected == 0 {
		result2 := tx.Create(&leaverecord)
		if result2.RowsAffected == 0 {
			tx.RollbackTo("savepoint")
			op.IsSuccess = false
			op.StatusCode = http.StatusNotFound
			op.Message = "Leaverecord for current year creation failed"
			return op
		}
	}
	println("Leave year", leaverecord.Year)

	_ = tx.Raw("select * from payroll.leave where empcode=?", input.Empcode).Find(&output)

	sumofcasual := 0
	for i := 0; i < len(output); i++ {
		sumofcasual = sumofcasual + output[i].Casual
	}

	fmt.Println("casual :", sumofcasual)
	_ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.leave").First(&input.Refno)

	result := db.Create(&input)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		op.IsSuccess = false
		op.StatusCode = 404
		op.Message = "Something went wrong"
		return op
	}
	tx.Commit()
	op.IsSuccess = true
	op.Message = "Success"
	op.StatusCode = http.StatusCreated
	op.Payload = input
	return op
}

func (leave *Leave) LeaveStatus(input model.Leave) dto.ResponseDto {
	var op dto.ResponseDto
	// Leave earns and consume
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output []model.Leave
	var output2 model.Leave
	// var leaveinput dto.LeaveCount
	date := time.Now()
	cdate := date.Format("2006-01-02")
	year := cdate[:4]
	firstdate := year + "-01-01"
	fmt.Println("Current date ", cdate)
	fmt.Println("First date ", firstdate)
	if input.Empcode <= 0 {
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Message = "Empcode is required"
		return op
	}
	result := db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ?", input.Empcode, firstdate, cdate).Find(&output)
	if result.RowsAffected == 0 {

		// output2.Casual = 0
		// output2.Medical = 0
		// output2.Earn = 0
		// output2.Sick = 0
		// output2.Maternity = 0
		// output2.Festival = 0
		// output2.Extraordi = 0
		// output2.Study = 0
		// output2.Leavewopay = 0
		// output2.Special = 0
		op.IsSuccess = false
		op.Message = "Employee not found"
		op.StatusCode = http.StatusNotFound
		op.Payload = nil
		return op
	} else {
		input.Casual = 0
		for i := 0; i < len(output); i++ {
			output2.Casual = output2.Casual + output[i].Casual
			output2.Earn = output2.Earn + output[i].Earn
			output2.Sick = output2.Sick + output[i].Sick
			output2.Medical = output2.Medical + output[i].Medical
			output2.Maternity = output2.Maternity + output[i].Maternity
			output2.Festival = output2.Festival + output[i].Festival
			output2.Extraordi = output2.Extraordi + output[i].Extraordi
			output2.Study = output2.Study + output[i].Study
			output2.Leavewopay = output2.Leavewopay + output[i].Leavewopay
			output2.Special = output2.Special + output[i].Special
		}
	}

	op.IsSuccess = true
	op.Message = "Success"
	op.StatusCode = 200
	op.Payload = output2
	return op
}

// result3 := db.Raw(`select a.empcode,a.empname ,a.joindate ,a.confirmdate
// from payroll.staffinformation a
// inner join payroll.salarystructure b
// on a.empcode =b.empcode
// where a.empcode=?`, input.Empcode).First(&output2)

// if result3.RowsAffected == 0 {
// 	op.IsSuccess = false
// 	op.Message = "Empcode not found"
// 	op.StatusCode = http.StatusNotFound
// 	op.Payload = nil
// 	return op
// }

// fmt.Println("Joining Date", output2.Joindate)
// // fmt.Println("Joining Date", output2.Joindate[:10])
// fmt.Println("confirm Date", output2.Confirmdate)
// // a := output2.Joindate
// a := output2.Joindate[:10]
// date := time.Now()
// fmt.Println(" Date - ", date)
// format := "2006-01-02"
// then, _ := time.Parse(format, a)

// fmt.Println("after time parse joining date : ", then)

// // b := output2.Confirmdate
// b := output2.Confirmdate[:10]
// now, _ := time.Parse(format, b)

// fmt.Println("after time parse confirmation date: ", now)

// diff := date.Sub(then)
// days := (int(diff.Hours() / 24))
// fmt.Println("days : ", days)
