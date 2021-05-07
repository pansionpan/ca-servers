package command

import (
	"errors"
	"fmt"
	"strconv"
)

type UserInput struct {
	Number      int
	SubjectName string
	FromDate    string
	ToDate      string
}

func GetUserInput() (*UserInput, error) {
	var fromDate, toDate, subjectName string
	var number int

	fmt.Printf(`请输入类型编号:
1:EnrollmentCredential
2:AuthorizationTicket
3:AuthorizationAuthority
4:EnrollmentAuthority
5:RootCa
6:CrlSigner
7:PseudonymTicket
8:PseudonymAuthority`)
	fmt.Scanln(&number)
	_, err := ChooseSubjectType(number)
	if err != nil {
		return nil, err
	}

	fmt.Printf("请输入名称: ")
	fmt.Scanln(&subjectName)
	if subjectName == "" {
		return nil, errors.New("the subjectName is wrong")
	}

	fmt.Printf("请输入起始年月日(例：20050101 最低年限为2005年1月1日): ")
	fmt.Scanln(&fromDate)
	if len(fromDate) != 8 {
		return nil, errors.New("输入的年月日有误")
	}

	days := fromDate[len(fromDate)-2:]
	months := fromDate[4:6]
	years := fromDate[:4]

	daysNum, err := strconv.Atoi(days)
	if err != nil {
		return nil, errors.New("日期输入有误")
	}
	monthsNum, err := strconv.Atoi(months)
	if err != nil {
		return nil, errors.New("月份输入有误")
	}

	yearsNum, err := strconv.Atoi(years)
	if err != nil {
		return nil, errors.New("年限输入有误")
	}

	if ( yearsNum % 4 == 0 && yearsNum % 100 != 0 ) || yearsNum % 400 ==0 {
		if monthsNum == 2 {
			if daysNum > 29 || daysNum < 01 {
				return nil, errors.New("闰年2月份日期输入有误")
			}
		}
	} else {
		if monthsNum == 2 {
			if daysNum > 28 || daysNum < 01 {
				return nil, errors.New("日期输入有误")
			}
		}
	}

	if yearsNum < 2005 {
		return nil, errors.New("年限输入有误")
	}
	if monthsNum > 12 || monthsNum < 01 {
		return nil, errors.New("月份输入有误")
	}
	if daysNum > 31 || daysNum < 01 {
		return nil, errors.New("日期输入有误")
	}

	fmt.Printf("请输入终止年月日(需大于起始时间): ")
	fmt.Scanln(&toDate)

	toDays := toDate[len(toDate)-2:]
	toMonths := toDate[4:6]
	toYears := toDate[:4]

	toDaysNum, err := strconv.Atoi(toDays)
	if err != nil {
		return nil, errors.New("日期输入有误")
	}
	toMonthsNum, err := strconv.Atoi(toMonths)
	if err != nil {
		return nil, errors.New("月份输入有误")
	}

	toYearsNum, err := strconv.Atoi(toYears)
	if err != nil {
		return nil, errors.New("年限输入有误")
	}

	if ( toYearsNum % 4 == 0 && toYearsNum % 100 != 0 ) || toYearsNum % 400 ==0 {
		if toMonthsNum == 2 {
			if toDaysNum > 29 || toDaysNum < 01 {
				return nil, errors.New("闰年2月份日期输入有误")
			}
		}
	} else {
		if toMonthsNum == 2 {
			if toDaysNum > 28 || toDaysNum < 01 {
				return nil, errors.New("日期输入有误")
			}
		}
	}

	if toYearsNum < 2005 {
		return nil, errors.New("年限输入有误")
	}
	if toMonthsNum > 12 || toMonthsNum < 01 {
		return nil, errors.New("月份输入有误")
	}
	if toDaysNum > 31 || toDaysNum < 01 {
		return nil, errors.New("日期输入有误")
	}


	return &UserInput{
		Number:      number,
		SubjectName: subjectName,
		FromDate:    fromDate,
		ToDate:      toDate,
	}, nil
}
