package excel

import (
	"github.com/tealeg/xlsx/v3"
	"pluto/usecase/dto"
	"time"
)

func ParseActivities(file *xlsx.File) dto.SetActivityRequest {
	request := new(dto.SetActivityRequest)

	for _, element := range activities(file) {
		request.Activities = append(request.Activities, element)
	}

	return *request
}

func activities(file *xlsx.File) []dto.Activity {
	request := *new([]dto.Activity)

	for _, sheet := range file.Sheets {
		for monthIdx := 2; monthIdx <= 2+27*12; monthIdx += 23 {
			monthRow, _ := sheet.Row(monthIdx)
			monthRawValue, e := monthRow.GetCell(5).GetTime(false)
			if e != nil {
				continue
			}
			monthValue := monthRawValue.Format("2006-01")

			lastRow, _ := sheet.Row(monthIdx + 2 + 4*5)
			var weekCount int
			if lastRow.GetCell(1).Value == "" {
				weekCount = 4
			} else {
				weekCount = 5
				monthIdx += 4
			}

			for weekIdx := monthIdx + 2; weekIdx <= monthIdx+2+4*weekCount; weekIdx += 4 {
				weekRow, _ := sheet.Row(weekIdx)
				secondFloorTeacherRow, _ := sheet.Row(weekIdx + 3)
				thirdFloorTeacherRow, _ := sheet.Row(weekIdx + 2)
				fourthFloorTeacherRow, _ := sheet.Row(weekIdx + 1)

				for dayIdx := 3; dayIdx <= 3+2*4; dayIdx += 2 {
					day := weekRow.GetCell(dayIdx).Value
					if len(day) == 0 {
						continue
					}
					if len(day) == 1 {
						day = "0" + day
					}
					date := monthValue + "-" + day

					t, _ := time.Parse("2006-01-02", date)
					weekday := t.Weekday().String()

					secondFloorTeacherName := secondFloorTeacherRow.GetCell(dayIdx + 1).Value
					thirdFloorTeacherName := thirdFloorTeacherRow.GetCell(dayIdx + 1).Value
					fourthFloorTeacherName := fourthFloorTeacherRow.GetCell(dayIdx + 1).Value

					request = append(request, dto.Activity{
						Date:                   date,
						Weekday:                weekday,
						SecondFloorTeacherName: secondFloorTeacherName,
						ThirdFloorTeacherName:  thirdFloorTeacherName,
						FourthFloorTeacherName: fourthFloorTeacherName,
					})
				}
			}
		}
	}

	return request
}
