package formPlatform

import (
	. "github.com/cdvelop/tinystring"
)

func (f *field) SetPropertiesFromInputTag(params ...string) {
	inputData := f.reflectStructField.Tag.Get("Input")
	if inputData != "" {
		// Get base params from input tag
		params = append(params, getParams(inputData))

		// Get input type from input tag
		if inputTypeFromLabel := Convert(inputData).Split("(")[0]; inputTypeFromLabel != "" {
			inputType = inputTypeFromLabel
		}
	}
	if f.customName == "" {
		f.customName = f.htmlName
	}

	// f.searchDataSourceImplementation(params...)

	options := f.separateOptions(params...)

	for _, option := range options {
		switch option {
		case "hidden":
			f.htmlName = option
		case "!required":
			f.allowSkipCompleted = true
		case `typing="hide"`:
			f.htmlName = "password"
		case "multiple":
			f.Multiple = option
		case "letters":
			f.Letters = true
		case "numbers":
			f.Numbers = true
		}

		switch {

		case Contains(option, "chars="):
			f.Characters = []rune(extractValue(option, "chars"))

		case Contains(option, "data="):
			extractData(extractValue(option, "data"), &f.DataSet)

		case Contains(option, "options="):
			extractData(extractValue(option, "options"), &f.options)

		case Contains(option, "class="):
			newClass := className(extractValue(option, "class"))
			exists := false
			for _, class := range f.Class {
				if class == newClass {
					exists = true
					break
				}
			}
			if !exists {
				f.Class = append(f.Class, newClass)
			}

		case Contains(option, "name="):
			f.Name, _ = Convert(option).KV("name")

		case Contains(option, "min="):
			f.Min, _ = Convert(option).KV("min")

		case Contains(option, "max="):
			f.Max, _ = Convert(option).KV("max")

		case Contains(option, "maxlength="):
			f.Maxlength, _ = Convert(option).KV("maxlength")

		case Contains(option, "placeholder="):
			f.PlaceHolder, _ = Convert(option).KV("placeholder")

		case Contains(option, "title="):
			f.Title, _ = Convert(option).KV("title")

		case Contains(option, "autocomplete="):
			f.Autocomplete, _ = Convert(option).KV("autocomplete")

		case Contains(option, "rows="):
			f.Rows, _ = Convert(option).KV("rows")

		case Contains(option, "cols="):
			f.Cols, _ = Convert(option).KV("cols")

		case Contains(option, "step="):
			f.Step, _ = Convert(option).KV("step")

		case Contains(option, "oninput="):
			f.Oninput, _ = Convert(option).KV("oninput")

		case Contains(option, "onkeyup="):
			f.Onkeyup, _ = Convert(option).KV("onkeyup")

		case Contains(option, "onchange="):
			f.Onchange, _ = Convert(option).KV("onchange")

		case Contains(option, "value="):
			f.Value, _ = Convert(option).KV("value")

		case Contains(option, "accept="):
			f.Accept, _ = Convert(option).KV("accept")

		}
	}

	if f.Name == "" {
		f.Name = f.customName
	}

	if f.Min != "" {
		f.Minimum, _ = Convert(f.Min).Int()
	}

	if f.Max != "" {
		f.Maximum, _ = Convert(f.Max).Int()
	}

	if f.htmlName != "hidden" {
		f.setDynamicTitle()
	} else {
		f.Title = ""
		f.PlaceHolder = ""
	}

	if len(f.options) == 0 {
		f.options = []map[string]string{
			{"": ""},
		}
	}
}
