package cli

import (
	"errors"
	"fmt"
	"testing"
)

func TestSuggestFlag(t *testing.T) {
	// Given
	app := testApp()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"a", "--another-flag"},
		{"hlp", "--help"},
		{"k", ""},
		{"s", "-s"},
	} {
		// When
		res := app.suggestFlag(app.Flags, testCase.provided)

		// Then
		expect(t, res, testCase.expected)
	}
}

func TestSuggestFlagHideHelp(t *testing.T) {
	// Given
	app := testApp()
	app.HideHelp = true

	// When
	res := app.suggestFlag(app.Flags, "hlp")

	// Then
	expect(t, res, "--fl")
}

func TestSuggestFlagFromError(t *testing.T) {
	// Given
	app := testApp()

	for _, testCase := range []struct {
		command, provided, expected string
	}{
		{"", "hel", "--help"},
		{"", "soccer", "--socket"},
		{"config", "anot", "--another-flag"},
	} {
		// When
		res, _ := app.suggestFlagFromError(
			errors.New(providedButNotDefinedErrMsg+testCase.provided),
			testCase.command,
		)

		// Then
		expect(t, res, fmt.Sprintf(didYouMeanTemplate+"\n\n", testCase.expected))
	}
}

func TestSuggestFlagFromErrorWrongError(t *testing.T) {
	// Given
	app := testApp()

	// When
	_, err := app.suggestFlagFromError(errors.New("invalid"), "")

	// Then
	expect(t, true, err != nil)
}

func TestSuggestFlagFromErrorWrongCommand(t *testing.T) {
	// Given
	app := testApp()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+"flag"),
		"invalid",
	)

	// Then
	expect(t, true, err != nil)
}

func TestSuggestFlagFromErrorNoSuggestion(t *testing.T) {
	// Given
	app := testApp()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+""),
		"",
	)

	// Then
	expect(t, true, err != nil)
}

func TestSuggestCommand(t *testing.T) {
	// Given
	app := testApp()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"conf", "config"},
		{"i", "i"},
		{"information", "info"},
		{"not-existing", "info"},
	} {
		// When
		res := suggestCommand(app.Commands, testCase.provided)

		// Then
		expect(t, res, fmt.Sprintf(didYouMeanTemplate, testCase.expected))
	}
}

func ExampleApp_Suggest() {
	app := &App{
		Name:                  "greet",
		Suggest:               true,
		HideHelp:              true,
		HideHelpCommand:       true,
		CustomAppHelpTemplate: "(this space intentionally left blank)\n",
		Flags: []Flag{
			&StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(cCtx *Context) error {
			fmt.Printf("Hello %v\n", cCtx.String("name"))
			return nil
		},
	}

	app.Run([]string{"greet", "--nema", "chipmunk"})
	// Output:
	// Incorrect Usage. flag provided but not defined: -nema
	//
	// Did you mean '--name'?
	//
	// (this space intentionally left blank)
}

func ExampleApp_Suggest_command() {
	app := &App{
		Name:                  "greet",
		Suggest:               true,
		HideHelp:              true,
		HideHelpCommand:       true,
		CustomAppHelpTemplate: "(this space intentionally left blank)\n",
		Flags: []Flag{
			&StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(cCtx *Context) error {
			fmt.Printf("Hello %v\n", cCtx.String("name"))
			return nil
		},
		Commands: []*Command{
			{
				Name:               "neighbors",
				CustomHelpTemplate: "(this space intentionally left blank)\n",
				Flags: []Flag{
					&BoolFlag{Name: "smiling"},
				},
				Action: func(cCtx *Context) error {
					if cCtx.Bool("smiling") {
						fmt.Println("😀")
					}
					fmt.Println("Hello, neighbors")
					return nil
				},
			},
		},
	}

	app.Run([]string{"greet", "neighbors", "--sliming"})
	// Output:
	// Incorrect Usage: flag provided but not defined: -sliming
	//
	// Did you mean '--smiling'?
	//
	// (this space intentionally left blank)
}
