package game_ui

func (ui *GameUI) drawBorders(draw drawFunc) {
	for x := range ui.cfg.XSize {
		draw(x, -1, '║')
		draw(x, ui.cfg.YSize, '║')
	}
	for y := range ui.cfg.YSize {
		draw(-1, y, '═')
		draw(ui.cfg.XSize, y, '═')
	}

	draw(-1, -1, '╔')
	draw(ui.cfg.XSize, -1, '╚')
	draw(-1, ui.cfg.YSize, '╗')
	draw(ui.cfg.XSize, ui.cfg.YSize, '╝')
}
