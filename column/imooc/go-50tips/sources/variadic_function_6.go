package main

import "fmt"

type record struct {
	name    string
	gender  string
	age     uint16
	city    string
	country string
}

func enroll(args ...interface{} /* name, gender, age, city = "Beijing", country = "China" */) (*record, error) {
	if len(args) > 5 || len(args) < 3 {
		return nil, fmt.Errorf("the number of arguments passed is wrong")
	}

	r := &record{
		city:    "Beijing", // 默认值：Beijing
		country: "China",   // 默认值：China
	}

	for i, v := range args {
		switch i {
		case 0: // name
			name, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("name is not passed as string")
			}
			r.name = name
		case 1: // gender
			gender, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("gender is not passed as string")
			}
			r.gender = gender
		case 2: // age
			age, ok := v.(int)
			if !ok {
				return nil, fmt.Errorf("age is not passed as int")
			}
			r.age = uint16(age)
		case 3: // city
			city, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("city is not passed as string")
			}
			r.city = city
		case 4: // country
			country, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("country is not passed as string")
			}
			r.country = country
		default:
			return nil, fmt.Errorf("unknown argument passed")
		}
	}

	return r, nil
}

func main() {
	r, _ := enroll("小明", "male", 23)
	fmt.Printf("%+v\n", *r)

	r, _ = enroll("小红", "female", 13, "Hangzhou")
	fmt.Printf("%+v\n", *r)

	r, _ = enroll("Leo Messi", "male", 33, "Barcelona", "Spain")
	fmt.Printf("%+v\n", *r)

	r, err := enroll("小吴", 21, "Suzhou")
	if err != nil {
		fmt.Println(err)
		return
	}
}
