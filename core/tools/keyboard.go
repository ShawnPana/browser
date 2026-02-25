package tools

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

var keyMap = map[string]input.Key{
	// Modifiers
	"control":  input.ControlLeft,
	"ctrl":     input.ControlLeft,
	"shift":    input.ShiftLeft,
	"alt":      input.AltLeft,
	"meta":     input.MetaLeft,
	"cmd":      input.MetaLeft,
	"command":  input.MetaLeft,

	// Special keys
	"enter":      input.Enter,
	"return":     input.Enter,
	"tab":        input.Tab,
	"escape":     input.Escape,
	"esc":        input.Escape,
	"space":      input.Space,
	"backspace":  input.Backspace,
	"delete":     input.Delete,
	"insert":     input.Insert,
	"home":       input.Home,
	"end":        input.End,
	"pageup":     input.PageUp,
	"pagedown":   input.PageDown,
	"arrowup":    input.ArrowUp,
	"arrowdown":  input.ArrowDown,
	"arrowleft":  input.ArrowLeft,
	"arrowright": input.ArrowRight,
	"capslock":   input.CapsLock,
	"numlock":    input.NumLock,

	// F-keys
	"f1":  input.F1,
	"f2":  input.F2,
	"f3":  input.F3,
	"f4":  input.F4,
	"f5":  input.F5,
	"f6":  input.F6,
	"f7":  input.F7,
	"f8":  input.F8,
	"f9":  input.F9,
	"f10": input.F10,
	"f11": input.F11,
	"f12": input.F12,

	// Letters
	"a": input.KeyA, "b": input.KeyB, "c": input.KeyC, "d": input.KeyD,
	"e": input.KeyE, "f": input.KeyF, "g": input.KeyG, "h": input.KeyH,
	"i": input.KeyI, "j": input.KeyJ, "k": input.KeyK, "l": input.KeyL,
	"m": input.KeyM, "n": input.KeyN, "o": input.KeyO, "p": input.KeyP,
	"q": input.KeyQ, "r": input.KeyR, "s": input.KeyS, "t": input.KeyT,
	"u": input.KeyU, "v": input.KeyV, "w": input.KeyW, "x": input.KeyX,
	"y": input.KeyY, "z": input.KeyZ,

	// Digits
	"0": input.Digit0, "1": input.Digit1, "2": input.Digit2, "3": input.Digit3,
	"4": input.Digit4, "5": input.Digit5, "6": input.Digit6, "7": input.Digit7,
	"8": input.Digit8, "9": input.Digit9,
}

var modifierKeys = map[string]bool{
	"control": true, "ctrl": true,
	"shift": true,
	"alt": true,
	"meta": true, "cmd": true, "command": true,
}

// Press sends a key combo like "Enter", "Control+a", "Shift+Enter".
func Press(page *rod.Page, combo string) error {
	parts := strings.Split(combo, "+")

	var modifiers []input.Key
	var keys []input.Key

	for _, part := range parts {
		name := strings.ToLower(strings.TrimSpace(part))
		if modifierKeys[name] {
			k, ok := keyMap[name]
			if !ok {
				return fmt.Errorf("unknown modifier: %s", part)
			}
			modifiers = append(modifiers, k)
		} else {
			k, ok := keyMap[name]
			if !ok {
				return fmt.Errorf("unknown key: %s", part)
			}
			keys = append(keys, k)
		}
	}

	ka := page.KeyActions()
	if len(modifiers) > 0 {
		ka = ka.Press(modifiers...)
	}
	if len(keys) > 0 {
		ka = ka.Type(keys...)
	}
	if len(modifiers) > 0 {
		ka = ka.Release(modifiers...)
	}
	return ka.Do()
}

// runeToKey maps a single rune to an input.Key for keyboard typing.
func runeToKey(r rune) (input.Key, bool) {
	if r >= 'a' && r <= 'z' {
		k, ok := keyMap[string(r)]
		return k, ok
	}
	if r >= 'A' && r <= 'Z' {
		k, ok := keyMap[string(unicode.ToLower(r))]
		return k, ok
	}
	if r >= '0' && r <= '9' {
		k, ok := keyMap[string(r)]
		return k, ok
	}
	switch r {
	case ' ':
		return input.Space, true
	case '\t':
		return input.Tab, true
	case '\n', '\r':
		return input.Enter, true
	}
	return 0, false
}

// KeyboardType types text character-by-character with real key events.
func KeyboardType(page *rod.Page, text string) error {
	for _, r := range text {
		k, ok := runeToKey(r)
		if ok {
			if r >= 'A' && r <= 'Z' {
				// Shift + key for uppercase
				if err := (page.KeyActions().Press(input.ShiftLeft).Type(k).Release(input.ShiftLeft).Do()); err != nil {
					return err
				}
			} else {
				if err := page.Keyboard.Type(k); err != nil {
					return err
				}
			}
		} else {
			// Fall back to InsertText for characters without a direct key mapping
			if err := page.InsertText(string(r)); err != nil {
				return err
			}
		}
	}
	return nil
}

// KeyboardInsertText inserts text directly without key events.
func KeyboardInsertText(page *rod.Page, text string) error {
	return page.InsertText(text)
}
