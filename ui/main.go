// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"recovery-tool/cmd"
)

const (
	defaultZipPath      = "../test/134_archive.zip"
	defaultVaultCount   = "1"
	defaultChain        = "60"
	defaultOutputPath   = "./output.yaml"
	defaultEciesPrivkey = "ea5db436b7508e5c8ec3ae17003bcb997c30e03c655f0dd2d1824ec93bd0501c"
	defaultMnemonic     = "amused garlic window please enrich sick gate ready owner giraffe elite umbrella hair seat punch seminar notable enroll wet asset outdoor inflict rich mushroom"
	defaultRsaPrivkey   = `-----BEGIN PRIVATE KEY-----
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
BkQYx9xJrcWT4gl+bVuezBX+nbcOyQ==
-----END PRIVATE KEY-----`
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Recovery Tool"),
			app.Size(unit.Dp(800), unit.Dp(800)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var okButton widget.Clickable

	var zipPathInput widget.Editor
	var vaultCountInput widget.Editor
	var chainInput widget.Editor
	var mnemonicInput widget.Editor
	var eciesPrivkeyInput widget.Editor
	var rsaPrivkeyInput widget.Editor
	var outputPathInput widget.Editor

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if okButton.Clicked(gtx) {
				zipPath := strings.TrimSpace(zipPathInput.Text())
				vaultCount, err := strconv.ParseInt(strings.TrimSpace(vaultCountInput.Text()), 10, 32)
				if err != nil {
					return fmt.Errorf("parse valut count error: %s", err.Error())
				}
				chain, err := strconv.ParseInt(strings.TrimSpace(chainInput.Text()), 10, 32)
				if err != nil {
					return fmt.Errorf("parse chain error: %s", err.Error())
				}
				mnemonic := strings.TrimSpace(mnemonicInput.Text())
				eciesPrivkey := strings.TrimSpace(eciesPrivkeyInput.Text())
				rsaPrivkey := strings.TrimSpace(rsaPrivkeyInput.Text())
				outputPath := strings.TrimSpace(outputPathInput.Text())

				params := cmd.RecoveryInput{
					ZipPath:      zipPath,
					UserMnemonic: mnemonic,
					EciesPrivKey: eciesPrivkey,
					RsaPrivKey:   rsaPrivkey,
					VaultCount:   int(vaultCount),
					CoinType:     []int{int(chain)},
				}

				result, err := cmd.RecoverKeys(params)
				if err != nil {
					return fmt.Errorf("derive keys error: %s", err.Error())
				}

				if err := cmd.SaveResult(&result, outputPath); err != nil {
					return fmt.Errorf("save result error: %s", err.Error())
				}

				fmt.Printf("Output the result to file `%s`\n", outputPath)
			}

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// zip path
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					zipPathInput.SingleLine = true
					zipPathInput.Alignment = text.Start
					ed := material.Editor(th, &zipPathInput, defaultZipPath)
					ed.Editor.SetText(defaultZipPath)
					margins := layout.Inset{
						Top:    unit.Dp(40),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// valut count
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					vaultCountInput.SingleLine = true
					vaultCountInput.Alignment = text.Start
					ed := material.Editor(th, &vaultCountInput, defaultVaultCount)
					ed.Editor.SetText(defaultVaultCount)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// chain
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					chainInput.SingleLine = true
					chainInput.Alignment = text.Start
					ed := material.Editor(th, &chainInput, defaultChain)
					ed.Editor.SetText(defaultChain)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// mnemonic
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					mnemonicInput.Alignment = text.Start
					ed := material.Editor(th, &mnemonicInput, defaultMnemonic)
					ed.Editor.SetText(defaultMnemonic)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// ecies private key
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					eciesPrivkeyInput.Alignment = text.Start
					ed := material.Editor(th, &eciesPrivkeyInput, defaultEciesPrivkey)
					ed.Editor.SetText(defaultEciesPrivkey)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// output path
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					outputPathInput.SingleLine = true
					outputPathInput.Alignment = text.Start
					ed := material.Editor(th, &outputPathInput, defaultOutputPath)
					ed.Editor.SetText(defaultOutputPath)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				// button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &okButton, "Recover Private Keys")
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Bottom: unit.Dp(25),
						Right:  unit.Dp(35),
						Left:   unit.Dp(35),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return btn.Layout(gtx)
						},
					)
				}),
				// empty spacer
				layout.Rigid(
					// The height of the spacer is 25 Device independent pixels
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
				// coincover private key
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rsaPrivkeyInput.Alignment = text.Start
					ed := material.Editor(th, &rsaPrivkeyInput, defaultRsaPrivkey)
					ed.Editor.SetText(defaultRsaPrivkey)
					margins := layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(170),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(170),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
