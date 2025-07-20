package internal

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// MAKE DYNAMIC
func PrintBanner() {
	ptermLogo, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("kaef", pterm.NewStyle(pterm.FgLightCyan)),
		putils.LettersFromStringWithStyle("fken", pterm.NewStyle(pterm.FgLightMagenta)),
	).Srender()

	pterm.DefaultCenter.Print("\n" + ptermLogo)
	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightCyan)).WithMargin(2).Sprint("[k2n] - ai based code/claim generation in your terminal or ci workflow"))

}

func PrintEnvTable(data map[string]string) error {
	var tableData pterm.TableData

	for k, v := range data {
		tableData = append(tableData, []string{
			pterm.White(k),
			pterm.LightMagenta(v),
		})
	}

	return pterm.DefaultTable.
		WithHasHeader(false).
		WithBoxed(false).
		WithSeparator("  ").
		WithData(tableData).
		Render()
}
