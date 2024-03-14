package webcontroller

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (wc *WebController) themeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(userStyleFromRequest(r)))
}

func userStyleFromRequest(r *http.Request) (s template.CSS) {
	// Get the chosen style from the URL
	var style = r.URL.Query().Get("style")
	var hue = -1

	if hueStr := r.URL.Query().Get("hue"); hueStr != "" {
		hue, _ = strconv.Atoi(hueStr)
	}

	// If the URL style was empty use the cookie value
	if style == "" {
		if cookie, err := r.Cookie("style"); err == nil {
			style = cookie.Value
		}
	}
	if hue == -1 {
		if cookie, err := r.Cookie("hue"); err == nil {
			hue, _ = strconv.Atoi(cookie.Value)
		}
	}

	return userStyle(style, hue)
}

func userStyle(style string, hue int) template.CSS {
	var (
		def      styleSheet
		light    styleSheet
		hasLight bool
	)

	switch style {
	default:
		fallthrough
	case "nord":
		def = nordDarkStyle
		light = nordLightStyle
		hasLight = true
	case "nord_dark":
		def = nordDarkStyle
	case "nord_light", "snowstorm":
		def = nordLightStyle
	case "solarized":
		def = solarizedDarkStyle
		light = solarizedLightStyle
		hasLight = true
	case "solarized_dark":
		def = solarizedDarkStyle
	case "solarized_light":
		def = solarizedLightStyle
	case "classic":
		def = classicStyle
		hue = -1 // Does not support custom hues
	case "purple_drain":
		def = purpleDrainStyle
		hue = -1 // Does not support custom hues
	case "maroon":
		def = maroonStyle
	case "hacker":
		def = hackerStyle
		hue = -1 // Does not support custom hues
	case "canta":
		def = cantaPixeldrainStyle
	case "skeuos":
		def = skeuosPixeldrainStyle
	case "sweet":
		def = sweetPixeldrainStyle
	case "adwaita":
		def = adwaitaDarkStyle
		light = adwaitaLightStyle
		hasLight = true
	case "adwaita_dark":
		def = adwaitaDarkStyle
	case "adwaita_light":
		def = adwaitaLightStyle
	}

	if hue >= 0 && hue <= 360 {
		def = def.withHue(hue)
		light = light.withHue(hue)
	}

	if hasLight {
		return template.CSS(def.withLight(light))
	} else {
		return template.CSS(def.String())
	}
}

type styleSheet struct {
	Link                HSL // Based on Highlight if undefined
	Input               CSS
	InputHover          CSS
	InputText           CSS
	InputDisabledText   CSS
	HighlightBackground CSS
	Highlight           HSL // Links, highlighted buttons, list navigation
	HighlightText       HSL // Text on buttons
	Danger              HSL
	ScrollbarForeground CSS // Based on Input if undefined
	ScrollbarHover      CSS // Based on ScrollbarForeground if undefined

	BackgroundColor   HSL
	Background        CSS
	BackgroundText    HSL
	BackgroundPattern CSS
	Navigation        CSS
	BodyColor         HSL
	BodyBackground    CSS
	BodyText          HSL
	Separator         HSL
	CardColor         HSL
	CardText          HSL

	// Colors to use in graphs
	Chart1 HSL
	Chart2 HSL
	Chart3 HSL
}

func (s styleSheet) withDefaults() styleSheet {
	// Set default colors
	var noColor = HSL{0, 0, 0}
	var defaultHSL = func(target *HSL, def HSL) {
		if *target == noColor {
			*target = def
		}
	}
	var defaultCSS = func(target *CSS, def CSS) {
		if *target == nil {
			*target = def
		}
	}
	defaultHSL(&s.Link, s.Highlight.Add(0, 0, -.05))
	defaultCSS(&s.ScrollbarForeground, s.Input)
	defaultCSS(&s.ScrollbarHover, s.InputHover)
	defaultHSL(&s.Chart1, s.Highlight)
	defaultHSL(&s.Chart2, s.Chart1.Add(120, 0, 0))
	defaultHSL(&s.Chart3, s.Chart2.Add(120, 0, 0))
	defaultCSS(&s.HighlightBackground, s.Highlight)
	defaultCSS(&s.Background, s.BackgroundColor)
	defaultCSS(&s.BackgroundPattern, s.BackgroundColor)
	defaultCSS(&s.Navigation, RawCSS("none"))
	defaultCSS(&s.BodyBackground, s.BodyColor)
	defaultHSL(&s.BackgroundText, s.BodyText)
	defaultHSL(&s.Separator, s.BodyColor.Add(0, 0, .05))

	return s
}

func (s styleSheet) withHue(hue int) styleSheet {
	s = s.withDefaults()

	var setHue = func(c CSS) CSS {
		if hsl, ok := c.(HSL); ok {
			hsl.Hue = hue
			return hsl
		} else {
			return c
		}
	}

	s.Input = setHue(s.Input)
	s.InputHover = setHue(s.InputHover)
	s.InputText = setHue(s.InputText)
	s.InputDisabledText = setHue(s.InputDisabledText)
	s.ScrollbarForeground = setHue(s.ScrollbarForeground)
	s.ScrollbarHover = setHue(s.ScrollbarHover)
	s.BackgroundColor.Hue = hue
	s.Background = setHue(s.Background)
	s.BackgroundText.Hue = hue
	s.BackgroundPattern = setHue(s.BackgroundPattern)
	s.Navigation = setHue(s.Navigation)
	s.BodyColor.Hue = hue
	s.BodyBackground = setHue(s.BodyBackground)
	s.BodyText.Hue = hue
	s.Separator.Hue = hue
	s.CardColor.Hue = hue
	s.CardText.Hue = hue
	return s
}

func (s styleSheet) String() string {
	s = s.withDefaults()

	return fmt.Sprintf(
		`:root {
	--link_color:                 %s;
	--input_background:           %s;
	--input_hover_background:     %s;
	--input_text:                 %s;
	--input_disabled_text:        %s;
	--highlight_background:       %s;
	--highlight_color:            %s;
	--highlight_text_color:       %s;
	--danger_color:               %s;
	--danger_text_color:          %s;
	--scrollbar_foreground_color: %s;
	--scrollbar_hover_color:      %s;

	--background_color:         %s;
	--background:               %s;
	--background_text_color:    %s;
	--background_pattern:       url("%s");
	--background_pattern_color: %s;
	--navigation_background:    %s;
	--body_color:               %s;
	--body_background:          %s;
	--body_text_color:          %s;
	--separator:                %s;
	--shaded_background:        %s;
	--card_color:               %s;

	--chart_1_color: %s;
	--chart_2_color: %s;
	--chart_3_color: %s;

	--shadow_color: %s;
}`,
		s.Link.CSS(),
		s.Input.CSS(),
		s.InputHover.CSS(),
		s.InputText.CSS(),
		s.InputDisabledText.CSS(),
		s.HighlightBackground.CSS(),
		s.Highlight.CSS(),
		s.HighlightText.CSS(),
		s.Danger.CSS(),
		s.HighlightText.CSS(),
		s.ScrollbarForeground.CSS(),
		s.ScrollbarHover.CSS(),
		s.BackgroundColor.CSS(),
		s.Background.CSS(),
		s.BackgroundText.CSS(),
		BackgroundTiles(),
		s.BackgroundPattern.CSS(),
		s.Navigation.CSS(),
		s.BodyColor.CSS(),
		s.BodyBackground.CSS(),
		s.BodyText.CSS(),
		s.Separator.CSS(),
		s.BodyColor.WithAlpha(0.8).CSS(), // shaded_background
		s.CardColor.CSS(),
		s.Chart1.CSS(),
		s.Chart2.CSS(),
		s.Chart3.CSS(),
		s.BodyColor.Darken(0.8).CSS(),
	)
}

func (dark styleSheet) withLight(light styleSheet) string {
	return fmt.Sprintf(
		`%s

@media (prefers-color-scheme: light) {
	%s
}`,
		dark.String(),
		light.String(),
	)
}

func BackgroundTiles() template.URL {
	var (
		now   = time.Now()
		month = now.Month()
		day   = now.Day()
		file  string
	)

	if now.Weekday() == time.Wednesday && rand.Intn(20) == 0 {
		file = "checker_wednesday"
	} else if month == time.August && day == 8 {
		file = "checker_dwarf"
	} else if month == time.August && day == 24 {
		file = "checker_developers"
	} else if month == time.October && day == 31 {
		file = "checker_halloween"
	} else if month == time.December && (day == 25 || day == 26 || day == 27) {
		file = "checker_christmas"
	} else {
		file = fmt.Sprintf("checker%d", now.UnixNano()%20)
	}

	return template.URL("/res/img/background_patterns/" + file + "_transparent.png")
}

// Following are all the available styles

var purpleDrainStyle = styleSheet{
	Input:               RGBA{77, 19, 236, .3},
	InputHover:          RGBA{77, 19, 236, .4},
	InputText:           HSL{0, 0, .9},
	InputDisabledText:   HSL{266, .85, .4},
	HighlightBackground: NewGradient(160, HSL{150, .84, .39}, HSL{85, .85, .35}),
	Highlight:           HSL{117, .63, .46},
	HighlightText:       HSL{0, 0, 0},
	Danger:              HSL{357, .63, .46},
	ScrollbarForeground: HSL{266, .85, .40},
	ScrollbarHover:      HSL{266, .85, .50},

	BackgroundColor:   HSL{270, .84, .12},
	Background:        NewGradient(120, HSL{255, .84, .16}, HSL{300, .85, .14}),
	BackgroundPattern: RawCSS("none"),
	Navigation:        RawCSS("none"),
	BodyColor:         HSL{274, .9, .14},
	BodyBackground:    NewGradient(120, HSL{255, .84, .18}, HSL{300, .85, .16}),
	BodyText:          HSL{0, 0, .8},
	CardColor:         HSL{275, .8, .18},
}

var classicStyle = styleSheet{
	Input:             HSL{0, 0, .18},
	InputHover:        HSL{0, 0, .22},
	InputText:         HSL{0, 0, .9},
	InputDisabledText: HSL{0, 0, .4},
	Highlight:         HSL{89, .60, .45},
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{339, .65, .31},

	BackgroundColor: HSL{0, 0, .08},
	BodyColor:       HSL{0, 0, .12},
	BodyText:        HSL{0, 0, .8},
	CardColor:       HSL{0, 0, .16},
}

var maroonStyle = styleSheet{
	Input:             HSL{0, .8, .20}, // hsl(0, 87%, 40%)
	InputHover:        HSL{0, .8, .25},
	InputText:         HSL{0, 0, 1},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{137, 1, .37}, //hsl(137, 100%, 37%)
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{9, .96, .42}, //hsl(9, 96%, 42%)

	BackgroundColor: HSL{0, .7, .05},
	BodyColor:       HSL{0, .8, .08}, // HSL{0, .8, .15},
	BodyText:        HSL{0, 0, .8},
	CardColor:       HSL{0, .9, .14},
}

var hackerStyle = styleSheet{
	Input:             HSL{0, 0, .1},
	InputHover:        HSL{0, 0, .15},
	InputText:         HSL{0, 0, 1},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{120, .8, .5},
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{0, 1, .4},

	BackgroundColor: HSL{0, 0, 0},
	BodyColor:       HSL{0, 0, .03},
	BodyText:        HSL{0, 0, .8},
	CardColor:       HSL{120, .4, .05},
}

var cantaPixeldrainStyle = styleSheet{
	Input:             HSL{167, .06, .30}, // hsl(167, 6%, 30%)
	InputHover:        HSL{167, .06, .35}, // hsl(167, 6%, 30%)
	InputText:         HSL{0, 0, 1},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{165, 1, .40}, // hsl(165, 100%, 40%)
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{40, 1, .5}, // hsl(40, 100%, 50%)

	BackgroundColor: HSL{180, .04, .16},
	BodyColor:       HSL{168, .05, .21},
	BodyText:        HSL{0, 0, .8},
	CardColor:       HSL{170, .05, .26},
}

var skeuosPixeldrainStyle = styleSheet{
	Input:             HSL{226, .15, .23}, //hsl(226, 15%, 23%)
	InputHover:        HSL{226, .15, .28},
	InputText:         HSL{60, .06, .93},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{282, .65, .54}, // hsl(282, 65%, 54%)
	HighlightText:     HSL{0, 0, 1},
	Danger:            HSL{0, .79, .43}, // hsl(0, 79%, 43%)

	BackgroundColor: HSL{232, .14, .11}, //hsl(232, 14%, 11%)
	BodyColor:       HSL{229, .14, .16}, // hsl(229, 14%, 16%)
	BodyText:        HSL{60, .06, .93},  // hsl(60, 6%, 93%)
	CardColor:       HSL{225, .14, .17}, // hsl(225, 14%, 17%)
}

var (
	nord0  = HSL{220, .16, .22} // Dark
	nord1  = HSL{222, .16, .28} // Dark
	nord2  = HSL{220, .17, .32} // Dark
	nord3  = HSL{220, .16, .36} // Light
	nord4  = HSL{219, .28, .88} // Light
	nord5  = HSL{218, .27, .92} // Light
	nord6  = HSL{218, .27, .94} // Light
	nord7  = HSL{179, .25, .65} // Teal
	nord8  = HSL{193, .43, .67} // Light blue
	nord11 = HSL{354, .42, .56} // Red
	nord14 = HSL{92, .28, .65}  // Green
)

var nordDarkStyle = styleSheet{
	Input:               nord3.Add(0, 0, 0.01),
	InputHover:          nord3.Add(0, 0, 0.03),
	InputText:           nord5,
	InputDisabledText:   nord0,
	Highlight:           nord14,
	HighlightText:       nord0,
	Danger:              nord11,
	ScrollbarForeground: nord7,
	ScrollbarHover:      nord8,

	BackgroundColor: nord0,
	BodyColor:       nord1,
	BodyText:        nord4,
	CardColor:       nord2,
}

var nordLightStyle = styleSheet{
	Link:                HSL{92, .40, .32},
	Input:               nord4.Add(0, 0, -0.04),
	InputHover:          nord4.Add(0, 0, -0.06),
	InputText:           nord0,
	InputDisabledText:   nord6,
	Highlight:           nord14,
	HighlightText:       nord1,
	Danger:              nord11,
	ScrollbarForeground: nord7,
	ScrollbarHover:      nord8,

	BackgroundColor:   nord4,
	BackgroundText:    nord0,
	BodyColor:         nord6,
	BodyText:          nord2,
	Separator:         nord4,
	BackgroundPattern: nord4,
	CardColor:         nord5,
}

var sweetPixeldrainStyle = styleSheet{
	Input:             HSL{229, .25, .18}, // hsl(229, 25%, 14%)
	InputHover:        HSL{229, .25, .23}, // hsl(229, 25%, 14%)
	InputText:         HSL{223, .13, .79},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{296, .88, .44},
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{356, 1, .64}, // hsl(356, 100%, 64%)

	BackgroundColor: HSL{225, .25, .06}, // hsl(225, 25%, 6%)
	BodyColor:       HSL{228, .25, .12}, // hsl(228, 25%, 12%)
	BodyText:        HSL{223, .13, .79}, // hsl(223, 13%, 79%)
	CardColor:       HSL{229, .25, .14}, // hsl(229, 25%, 14%)
}

var adwaitaDarkStyle = styleSheet{
	Input:             RGBA{255, 255, 255, .06},
	InputHover:        RGBA{255, 255, 255, .11},
	InputText:         HSL{0, 0, 1},
	InputDisabledText: HSL{0, 0, .5},
	Highlight:         HSL{152, .62, .39}, // hsl(152, 62%, 39%)
	HighlightText:     HSL{0, 0, 0},
	Danger:            HSL{9, 1, .69}, // hsl(9, 100%, 69%)

	BackgroundColor: HSL{0, 0, .12},
	BodyColor:       HSL{0, 0, .14},
	BodyText:        HSL{0, 0, 1},
	CardColor:       HSL{0, 0, .21},
}

var adwaitaLightStyle = styleSheet{
	Input:             RGBA{0, 0, 0, .06},
	InputHover:        RGBA{0, 0, 0, .11},
	InputText:         HSL{0, 0, .2},
	InputDisabledText: HSL{0, 0, .7},
	Highlight:         HSL{152, .62, .47}, // hsl(152, 62%, 47%)
	HighlightText:     HSL{0, 0, 1},
	Danger:            HSL{356, .75, .43}, // hsl(356, 75%, 43%)

	BackgroundColor: HSL{0, 0, .92},
	BodyColor:       HSL{0, 0, .98},
	BodyText:        HSL{0, 0, .2},
	Separator:       HSL{0, 0, .86},
	CardColor:       HSL{0, 0, 1},
}

var solarizedDarkStyle = styleSheet{
	Input:             HSL{192, .81, .18}, // hsl(194, 14%, 40%)
	InputHover:        HSL{192, .81, .23}, // hsl(196, 13%, 45%)
	InputText:         HSL{180, .07, .80}, // hsl(44, 87%, 94%)
	InputDisabledText: HSL{194, .14, .30}, // hsl(194, 14%, 40%)
	Highlight:         HSL{68, 1, .30},    // hsl(68, 100%, 30%)
	HighlightText:     HSL{192, .81, .14}, // hsl(192, 100%, 11%)
	Danger:            HSL{1, .71, .52},   // hsl(1, 71%, 52%)

	BackgroundColor: HSL{192, 1, .11},   //hsl(192, 100%, 11%)
	BodyColor:       HSL{192, .81, .14}, // hsl(192, 81%, 14%)
	BodyText:        HSL{180, .07, .60}, // hsl(180, 7%, 60%)
	CardColor:       HSL{192, .81, .16},
}

var solarizedLightStyle = styleSheet{
	Input:             HSL{46, .42, .84},
	InputHover:        HSL{46, .42, .80},
	InputText:         HSL{194, .14, .20}, // hsl(192, 81%, 14%)
	InputDisabledText: HSL{44, .87, .94},
	Highlight:         HSL{68, 1, .30}, // hsl(68, 100%, 30%)
	HighlightText:     HSL{44, .87, .94},
	Danger:            HSL{1, .71, .52}, // hsl(1, 71%, 52%)

	BackgroundColor: HSL{46, .42, .88}, // hsl(46, 42%, 88%)
	BackgroundText:  HSL{192, .81, .14},
	BodyColor:       HSL{44, .87, .94},  // hsl(44, 87%, 94%)
	BodyText:        HSL{194, .14, .40}, // hsl(194, 14%, 40%)
	Separator:       HSL{46, .42, .88},
	CardColor:       HSL{44, .87, .92},
}
