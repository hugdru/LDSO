package data

import (
	"fmt"
	"testing"
)

func Test_1(t *testing.T) {
	sub1 := Sub_Criteria{"test1", 14}
	sub2 := Sub_Criteria{"test2", 20}

	var subs_array []Sub_Criteria = make([]Sub_Criteria, 2)
	subs_array[0] = sub1
	subs_array[1] = sub2

	crit := Criteria_Set{Criteria{"test", 10}, subs_array}

	fmt.Println(sub1)
	fmt.Println(sub2)
	fmt.Println(crit)
}

func Test_2(t *testing.T) {
	var var1 Sub_Criteria = Sub_Criteria{"ola", 23}
	aVar := make([]Sub_Criteria, 1)
	aVar[0] = var1

	var test_var Criteria_Set = Criteria_Set{Criteria: Criteria{Weight: 30,
			Name: "Access"}, Sub_Criterias: aVar}
	fmt.Println(test_var)

	test_var.SetSub_Criterias(Sub_Criteria{"Wall", 10},
			Sub_Criteria{"Ceil", 5}, Sub_Criteria{"Floor", 15})
	fmt.Println(test_var)

	test_var.Criteria.Weight = 35
	test_var.AppendSub_Criteria(Sub_Criteria{"Glass", 5})
	fmt.Println(test_var)
}
