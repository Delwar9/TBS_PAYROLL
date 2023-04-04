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

type Leavecheck struct{}

func (leavecheck Leavecheck) CheckEarnLeave(input model.Leaverecord) dto.ResponseDto {
	var op dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	tx := db.Begin()
	tx.SavePoint("savePoint")
	var yearoutput model.Leaverecord
	var output []model.Leaverecord
	var leaverecordop model.Leaverecord
	var idinput model.Leaverecord
	var weekendinput model.Weekend
	var weekendop model.Weekend
	var holidayinfoinput model.Holidayinfo
	var holidayop model.Holidayinfo
	var output3 []dto.LeaveCount

	dt := time.Now()
	currentyear := dt.Year()
	fmt.Println("Year - - - ", currentyear)
	result1 := db.Raw("select * from payroll.weekend where year=?", currentyear).First(&weekendinput)
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savePoint")
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Payload = nil
		op.Message = "Please Update weekend list of this year"
		return op
	}
	result2 := db.Raw("select * from payroll.holidayinfo where year=?", currentyear).First(&holidayinfoinput)
	if result2.RowsAffected == 0 {
		tx.RollbackTo("savePoint")
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Payload = nil
		op.Message = "Please Update Holiday Info list of this year"
		return op
	}

	result3 := db.Raw("select * from payroll.leaverecord where year=?", currentyear).First(&holidayinfoinput)
	if result3.RowsAffected > 0 {
		tx.RollbackTo("savePoint")
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Payload = nil
		op.Message = "Leave reord of this year is already updated!"
		return op
	}

	_ = db.Raw("select max(year) as year from payroll.leaverecord").First(&yearoutput)
	fmt.Println(yearoutput.Year)
	//result := tx.Raw(`select * from payroll.leaverecord where year = (select max(year) from payroll.leaverecord) order by empcode`).Find(&output)
	result := tx.Raw(`select * from payroll.leaverecord where year = ? order by empcode`, yearoutput.Year).Find(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint")
		op.IsSuccess = false
		op.StatusCode = http.StatusNotFound
		op.Payload = nil
		op.Message = "No leave record found"
		return op
	}

	if (currentyear > yearoutput.Year) && (currentyear == yearoutput.Year+1) {
		fmt.Println("inside if condition")
		_ = db.Raw("select coalesce ((max(id)+1),1) from payroll.leaverecord").First(&idinput.Id)

		_ = db.Raw("select count(noofdays) noofdays from payroll.holidayinfo where year=?", currentyear).First(&holidayop)
		_ = db.Raw("select count(noofweek) noofweek from payroll.weekend where year=?", currentyear).First(&weekendop)
		for i := 0; i < len(output); i++ {
			fmt.Println("inside loop")

			_ = db.Raw("select * from payroll.leaverecord where empcode=? and year=?", output[i].Empcode, yearoutput.Year).First(&leaverecordop)
			fmt.Println("Current Date", currentyear)

			firstday := strconv.Itoa(currentyear) + "-01-01"
			lastday := strconv.Itoa(currentyear) + "-12-31"

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

			_ = db.Raw("select * from payroll.leave where empcode=? and leavedate between ? and ? ", input.Empcode, firstday, lastday).Find(&output3)

			sumofearnleave := 0
			sumofallleave := 0
			for i := 0; i < len(output3); i++ {

				sumofearnleave = sumofearnleave + output3[i].Earn
				sumofallleave = (sumofallleave + output3[i].Casual + output3[i].Medical + output3[i].Special + output3[i].Earn + output3[i].Study +
					output3[i].Maternity + output3[i].Leavewopay + output3[i].Extraordi + output3[i].Festival + output3[i].Sick)

			}

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

			fmt.Println("Max ID", idinput.Id)
			output[i].Id = idinput.Id
			output[i].Year = currentyear
			//output[i].Empcode = output[i].Empcode
			output[i].Earnleave = leaverecordop.El_earned_nextyear + leaverecordop.Rest_of_leave
			output[i].Earnconsume = 0
			output[i].Rest_of_leave = output[i].Earnleave
			output[i].El_earned_nextyear = actual_earnleave

			result2 := tx.Create(&output[i])
			if result2.Error != nil {
				tx.RollbackTo("savePoint")
				op.IsSuccess = false
				op.StatusCode = http.StatusInternalServerError
				op.Payload = nil
				op.Message = "Error while creating leave record"
				return op
			}
			idinput.Id = output[i].Id + 1
		}

	}

	tx.Commit()
	op.IsSuccess = true
	op.StatusCode = http.StatusAccepted
	op.Message = "Successfully Updated"
	op.Payload = output
	return op
}
