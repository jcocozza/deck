package draw

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"os"
)

type Text struct {
	C color.Color
	Size int
}

type Theme struct {
	Background     Text // lt1
	Foreground     Text // dk1
	Heading        Text // dk2 or accent1
	SubHeading     Text // accent2
	Text           Text // dk1
	CardBackground Text // lt2
	Link           Text // hlink
	VisitedLink    Text // folHlink
	Primary        Text // accent1
	Secondary      Text // accent3 or accent5
}

var DefaultTheme = Theme{
	Background:     Text{ Size: 14, C: color.RGBA{R: 255, G: 255, B: 255, A: 255} }, // White (lt1)
	Foreground:     Text{ Size: 14, C: color.RGBA{R: 0, G: 0, B: 0, A: 255} },       // Black (dk1)
	Heading:        Text{ Size: 18, C: color.RGBA{R: 68, G: 114, B: 196, A: 255} },  // Accent1 (blue)
	SubHeading:     Text{ Size: 16, C: color.RGBA{R: 237, G: 125, B: 49, A: 255} },  // Accent2 (orange)
	Text:           Text{ Size: 14, C: color.RGBA{R: 0, G: 0, B: 0, A: 255} },       // Same as Foreground
	CardBackground: Text{ Size: 14, C: color.RGBA{R: 242, G: 242, B: 242, A: 255} }, // Light gray (lt2)
	Link:           Text{ Size: 14, C: color.RGBA{R: 5, G: 99, B: 193, A: 255} },    // Hyperlink blue
	VisitedLink:    Text{ Size: 14, C: color.RGBA{R: 149, G: 79, B: 114, A: 255} },  // Visited purple
	Primary:        Text{ Size: 14, C: color.RGBA{R: 68, G: 114, B: 196, A: 255} },  // Accent1 again
	Secondary:      Text{ Size: 14, C: color.RGBA{R: 165, G: 165, B: 165, A: 255} }, // Neutral gray (accent5-like)
}


func readTheme(path string) (*ThemeXML, error) {
	fbytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t ThemeXML
	err = xml.Unmarshal(fbytes, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ReadTheme(path string) (Theme, error) {
	theme, err := readTheme(path)
	if err != nil {
		return Theme{}, err
	}
	return ThemeXMLToTheme(theme)
}

func HexToRGBA(hex string) (color.Color, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func ThemeXMLToTheme(t *ThemeXML) (Theme, error) {
	background, err := HexToRGBA(t.ThemeElements.ClrScheme.Lt1.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	foreground, err := HexToRGBA(t.ThemeElements.ClrScheme.Dk1.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	heading, err := HexToRGBA(t.ThemeElements.ClrScheme.Dk2.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	subheading, err := HexToRGBA(t.ThemeElements.ClrScheme.Accent2.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	text, err := HexToRGBA(t.ThemeElements.ClrScheme.Dk1.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	cardbackground, err := HexToRGBA(t.ThemeElements.ClrScheme.Lt2.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	link, err := HexToRGBA(t.ThemeElements.ClrScheme.Hlink.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	visitedlink, err := HexToRGBA(t.ThemeElements.ClrScheme.FolHlink.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	primary, err := HexToRGBA(t.ThemeElements.ClrScheme.Accent1.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	secondary, err := HexToRGBA(t.ThemeElements.ClrScheme.Accent3.SrgbClr.Val)
	if err != nil {
		return Theme{}, err
	}
	return Theme{
		Background:     Text{ Size: 14, C: background },
		Foreground:     Text{ Size: 14, C: foreground },
		Heading:        Text{ Size: 18, C: heading },
		SubHeading:     Text{ Size: 16, C: subheading },
		Text:           Text{ Size: 14, C: text },
		CardBackground: Text{ Size: 14, C: cardbackground },
		Link:           Text{ Size: 14, C: link },
		VisitedLink:    Text{ Size: 14, C: visitedlink },
		Primary:        Text{ Size: 14, C: primary },
		Secondary:      Text{ Size: 14, C: secondary },
	}, nil
}

// Theme was generated 2025-05-25 09:16:07 by https://xml-to-go.github.io/ in Ukraine.
type ThemeXML struct {
	XMLName       xml.Name `xml:"theme"`
	Text          string   `xml:",chardata"`
	A             string   `xml:"a,attr"`
	Name          string   `xml:"name,attr"`
	ThemeElements struct {
		Text      string `xml:",chardata"`
		ClrScheme struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			Dk1  struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"dk1"`
			Lt1 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"lt1"`
			Dk2 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"dk2"`
			Lt2 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"lt2"`
			Accent1 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent1"`
			Accent2 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent2"`
			Accent3 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent3"`
			Accent4 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent4"`
			Accent5 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent5"`
			Accent6 struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"accent6"`
			Hlink struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"hlink"`
			FolHlink struct {
				Text    string `xml:",chardata"`
				SrgbClr struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"folHlink"`
		} `xml:"clrScheme"`
		FontScheme struct {
			Text      string `xml:",chardata"`
			Name      string `xml:"name,attr"`
			MajorFont struct {
				Text  string `xml:",chardata"`
				Latin struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"latin"`
				Ea struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"ea"`
				Cs struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"cs"`
			} `xml:"majorFont"`
			MinorFont struct {
				Text  string `xml:",chardata"`
				Latin struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"latin"`
				Ea struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"ea"`
				Cs struct {
					Text     string `xml:",chardata"`
					Typeface string `xml:"typeface,attr"`
				} `xml:"cs"`
			} `xml:"minorFont"`
		} `xml:"fontScheme"`
		FmtScheme struct {
			Text         string `xml:",chardata"`
			Name         string `xml:"name,attr"`
			FillStyleLst struct {
				Text      string `xml:",chardata"`
				SolidFill []struct {
					Text      string `xml:",chardata"`
					SchemeClr struct {
						Text string `xml:",chardata"`
						Val  string `xml:"val,attr"`
						Tint struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"tint"`
						Shade struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"shade"`
					} `xml:"schemeClr"`
				} `xml:"solidFill"`
			} `xml:"fillStyleLst"`
			LnStyleLst struct {
				Text string `xml:",chardata"`
				Ln   []struct {
					Text      string `xml:",chardata"`
					W         string `xml:"w,attr"`
					Cap       string `xml:"cap,attr"`
					Cmpd      string `xml:"cmpd,attr"`
					Algn      string `xml:"algn,attr"`
					SolidFill struct {
						Text      string `xml:",chardata"`
						SchemeClr struct {
							Text  string `xml:",chardata"`
							Val   string `xml:"val,attr"`
							Shade struct {
								Text string `xml:",chardata"`
								Val  string `xml:"val,attr"`
							} `xml:"shade"`
						} `xml:"schemeClr"`
					} `xml:"solidFill"`
					PrstDash struct {
						Text string `xml:",chardata"`
						Val  string `xml:"val,attr"`
					} `xml:"prstDash"`
					NoFill string `xml:"noFill"`
				} `xml:"ln"`
			} `xml:"lnStyleLst"`
			EffectStyleLst struct {
				Text        string `xml:",chardata"`
				EffectStyle []struct {
					Text      string `xml:",chardata"`
					EffectLst struct {
						Text        string `xml:",chardata"`
						Blur        string `xml:"blur"`
						FillOverlay struct {
							Text      string `xml:",chardata"`
							Blend     string `xml:"blend,attr"`
							SolidFill struct {
								Text      string `xml:",chardata"`
								SchemeClr struct {
									Text  string `xml:",chardata"`
									Val   string `xml:"val,attr"`
									Shade struct {
										Text string `xml:",chardata"`
										Val  string `xml:"val,attr"`
									} `xml:"shade"`
								} `xml:"schemeClr"`
							} `xml:"solidFill"`
						} `xml:"fillOverlay"`
					} `xml:"effectLst"`
				} `xml:"effectStyle"`
			} `xml:"effectStyleLst"`
			BgFillStyleLst struct {
				Text      string `xml:",chardata"`
				SolidFill []struct {
					Text      string `xml:",chardata"`
					SchemeClr struct {
						Text string `xml:",chardata"`
						Val  string `xml:"val,attr"`
						Tint struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"tint"`
						SatMod struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"satMod"`
					} `xml:"schemeClr"`
				} `xml:"solidFill"`
				GradFill struct {
					Text         string `xml:",chardata"`
					RotWithShape string `xml:"rotWithShape,attr"`
					GsLst        struct {
						Text string `xml:",chardata"`
						Gs   []struct {
							Text      string `xml:",chardata"`
							Pos       string `xml:"pos,attr"`
							SchemeClr struct {
								Text string `xml:",chardata"`
								Val  string `xml:"val,attr"`
								Tint struct {
									Text string `xml:",chardata"`
									Val  string `xml:"val,attr"`
								} `xml:"tint"`
								SatMod struct {
									Text string `xml:",chardata"`
									Val  string `xml:"val,attr"`
								} `xml:"satMod"`
								Shade struct {
									Text string `xml:",chardata"`
									Val  string `xml:"val,attr"`
								} `xml:"shade"`
								LumMod struct {
									Text string `xml:",chardata"`
									Val  string `xml:"val,attr"`
								} `xml:"lumMod"`
							} `xml:"schemeClr"`
						} `xml:"gs"`
					} `xml:"gsLst"`
					Lin struct {
						Text   string `xml:",chardata"`
						Ang    string `xml:"ang,attr"`
						Scaled string `xml:"scaled,attr"`
					} `xml:"lin"`
				} `xml:"gradFill"`
			} `xml:"bgFillStyleLst"`
		} `xml:"fmtScheme"`
	} `xml:"themeElements"`
	ObjectDefaults    string `xml:"objectDefaults"`
	ExtraClrSchemeLst string `xml:"extraClrSchemeLst"`
	CustClrLst        struct {
		Text    string `xml:",chardata"`
		CustClr []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name,attr"`
			SrgbClr struct {
				Text string `xml:",chardata"`
				Val  string `xml:"val,attr"`
			} `xml:"srgbClr"`
		} `xml:"custClr"`
	} `xml:"custClrLst"`
	ExtLst struct {
		Text string `xml:",chardata"`
		Ext  struct {
			Text        string `xml:",chardata"`
			URI         string `xml:"uri,attr"`
			ThemeFamily struct {
				Text  string `xml:",chardata"`
				Thm15 string `xml:"thm15,attr"`
				Name  string `xml:"name,attr"`
				ID    string `xml:"id,attr"`
				Vid   string `xml:"vid,attr"`
			} `xml:"themeFamily"`
		} `xml:"ext"`
	} `xml:"extLst"`
}
