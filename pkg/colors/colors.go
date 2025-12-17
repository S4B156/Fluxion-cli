package colors

import "github.com/fatih/color"

var (
    // --- Standard Statuses  ---
    Success  = color.New(color.FgGreen, color.Bold) 
    Error    = color.New(color.FgRed, color.Bold)   
    Warn     = color.New(color.FgYellow, color.Bold)
    Info     = color.New(color.FgCyan)              
    Debug    = color.New(color.FgWhite, color.Faint)
    Critical = color.New(color.BgRed, color.FgWhite, color.Bold)

    // --- GhostDraft Specifics ---
    Detected = color.New(color.FgMagenta, color.Bold) // Для "Проект найден"
    Path     = color.New(color.FgWhite, color.Bold)   // Для выделения путей и имен файлов

    // --- Technology Brand Colors (Ваши и Новые) ---
	Docker   = color.RGB(28, 144, 237)    // Синий
	Spring   = color.RGB(109, 179, 63)    // Весенний Зеленый
)

func Init() {
    // color.NoColor = true // для отключения в CI
}