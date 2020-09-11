package attendee

import (
	"encoding/xml"
	"strconv"
)

type Attendee struct {
	Name  string
	Age   int
	Phone string
}

func (a *Attendee) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{}

	tokens = append(tokens, xml.StartElement{
		Name: xml.Name{"", "attendee"}})

	tokens = append(tokens, xml.StartElement{Name: xml.Name{"", "name"}})
	tokens = append(tokens, xml.CharData(a.Name))
	tokens = append(tokens, xml.EndElement{Name: xml.Name{"", "name"}})

	tokens = append(tokens, xml.StartElement{Name: xml.Name{"", "age"}})
	tokens = append(tokens, xml.CharData(strconv.Itoa(a.Age)))
	tokens = append(tokens, xml.EndElement{Name: xml.Name{"", "age"}})

	tokens = append(tokens, xml.StartElement{Name: xml.Name{"", "phone"}})
	tokens = append(tokens, xml.CharData(a.Phone))
	tokens = append(tokens, xml.EndElement{Name: xml.Name{"", "phone"}})

	tokens = append(tokens, xml.StartElement{Name: xml.Name{"", "website"}})
	tokens = append(tokens, xml.CharData("https://www.gophercon.com/speaker/"+a.Name))
	tokens = append(tokens, xml.EndElement{Name: xml.Name{"", "website"}})

	tokens = append(tokens, xml.EndElement{Name: xml.Name{"", "attendee"}})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}
