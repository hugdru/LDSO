package data

import (
	"fmt"
	"testing"
)

func Test_Evaluation(t *testing.T) {
	var test_var1 Criteria_Set = Criteria_Set{Criteria{"Access", 30}, nil}
	test_var1.SetSub_Criterias(Sub_Criteria{"Wall", 30},
			Sub_Criteria{"Ceil", 10})
	var test_var2 Criteria_Set = Criteria_Set{Criteria{"Other", 70}, nil}
	test_var2.SetSub_Criterias(Sub_Criteria{"W", 60},
			Sub_Criteria{"C", 10})

	criterias := make([]Criteria_Set, 0)
	criterias = append(criterias, test_var1)
	criterias = append(criterias, test_var2)

	evaluation := Evaluation{Property{"Rua", Owner{"jose"}, "image.jpg"},
			criterias, nil, 0}
	evaluation.Evaluate()

	fmt.Println(evaluation)
	fmt.Println(evaluation.Value)
}

