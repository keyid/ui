package main

//go:generate mkembeddedfs --no-modtime --output images/fs_gen.go --pkg images --name FS --strip images/images images/images

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ui"
	"github.com/richardwilkes/ui/Demo/images"
	"github.com/richardwilkes/ui/app"
	"github.com/richardwilkes/ui/border"
	"github.com/richardwilkes/ui/cursor"
	"github.com/richardwilkes/ui/draw"
	"github.com/richardwilkes/ui/draw/align"
	"github.com/richardwilkes/ui/event"
	"github.com/richardwilkes/ui/font"
	"github.com/richardwilkes/ui/keys"
	"github.com/richardwilkes/ui/layout"
	"github.com/richardwilkes/ui/layout/flex"
	"github.com/richardwilkes/ui/layout/flow"
	"github.com/richardwilkes/ui/menu"
	"github.com/richardwilkes/ui/menu/appmenu"
	"github.com/richardwilkes/ui/menu/editmenu"
	"github.com/richardwilkes/ui/menu/filemenu"
	"github.com/richardwilkes/ui/menu/helpmenu"
	"github.com/richardwilkes/ui/menu/windowmenu"
	"github.com/richardwilkes/ui/widget"
	"github.com/richardwilkes/ui/widget/button"
	"github.com/richardwilkes/ui/widget/checkbox"
	"github.com/richardwilkes/ui/widget/imagebutton"
	"github.com/richardwilkes/ui/widget/imagelabel"
	"github.com/richardwilkes/ui/widget/label"
	"github.com/richardwilkes/ui/widget/list"
	"github.com/richardwilkes/ui/widget/popupmenu"
	"github.com/richardwilkes/ui/widget/radiobutton"
	"github.com/richardwilkes/ui/widget/scrollarea"
	"github.com/richardwilkes/ui/widget/separator"
	"github.com/richardwilkes/ui/widget/textfield"
	"github.com/richardwilkes/ui/widget/tooltip"
	"github.com/richardwilkes/ui/widget/webview"
	"github.com/richardwilkes/ui/window"
)

var (
	aboutWindow ui.Window
	imagesFS    = images.FS.FileSystem(("images/images"))
)

func main() {
	// event.TraceLogger = &jot.Logger{}
	// event.TraceAllEvents = true
	// event.TraceEventTypes = append(event.TraceEventTypes, event.MouseDownType, event.MouseDraggedType, event.MouseUpType)
	handlers := app.EventHandlers()
	handlers.Add(event.AppPopulateMenuBarType, createMenus)
	handlers.Add(event.AppWillFinishStartupType, finishStartup)
	app.Start()
}

func finishStartup(_ event.Event) {
	w1 := createButtonsWindow("Demo #1", geom.Point{})
	frame1 := w1.Frame()
	createButtonsWindow("Demo #2", geom.Point{X: frame1.X + frame1.Width + 5, Y: frame1.Y})
}

func createMenus(evt event.Event) {
	bar := menu.AppBar(evt.(*event.AppPopulateMenuBar).ID())
	_, aboutItem, prefsItem := appmenu.Install(bar)
	aboutItem.EventHandlers().Add(event.SelectionType, createAboutWindow)
	prefsItem.EventHandlers().Add(event.SelectionType, createPreferencesWindow)
	bar.AppendMenu(newFileMenu())
	editmenu.Install(bar)
	windowmenu.Install(bar)
	helpmenu.Install(bar)
}

func newFileMenu() menu.Menu {
	fileMenu := menu.NewMenu("File")
	fileMenu.AppendItem(menu.NewItemWithKey("Open", keys.VirtualKeyO, nil))
	fileMenu.AppendItem(menu.NewSeparator())
	fileMenu.AppendItem(filemenu.NewCloseKeyWindowItem())
	return fileMenu
}

func createButtonsWindow(title string, where geom.Point) ui.Window {
	wnd := window.NewWindow(where, window.StdWindowMask)
	wnd.SetTitle(title)

	content := wnd.Content()
	content.SetBorder(border.NewEmpty(geom.NewUniformInsets(10)))
	lay := flex.NewLayout(content)
	lay.VSpacing = 10

	buttonsPanel := createButtonsPanel()
	flexData := flex.NewData()
	flexData.HGrab = true
	buttonsPanel.SetLayoutData(flexData)
	content.AddChild(buttonsPanel)

	addSeparator(content)

	checkBoxPanel := createCheckBoxPanel()
	checkBoxPanel.SetLayoutData(flexData.Clone())
	content.AddChild(checkBoxPanel)

	addSeparator(content)

	radioButtonsPanel := createRadioButtonsPanel()
	radioButtonsPanel.SetLayoutData(flexData.Clone())
	content.AddChild(radioButtonsPanel)

	addSeparator(content)

	popupMenusPanel := createPopupMenusPanel()
	popupMenusPanel.SetLayoutData(flexData.Clone())
	content.AddChild(popupMenusPanel)

	addSeparator(content)

	wrapper := widget.NewBlock()
	lay = flex.NewLayout(wrapper)
	lay.Columns = 2
	lay.EqualColumns = true
	lay.HSpacing = 10
	flexData = flexData.Clone()
	flexData.HAlign = align.Fill
	wrapper.SetLayoutData(flexData)
	textFieldsPanel := createTextFieldsPanel()
	textFieldsPanel.SetLayoutData(flexData.Clone())
	wrapper.AddChild(textFieldsPanel)
	wrapper.AddChild(createListPanel())
	content.AddChild(wrapper)

	addSeparator(content)

	if title == "Demo #1" {
		wv := webview.NewWebView(wnd)
		flexData = flex.NewData()
		flexData.HAlign = align.Fill
		flexData.VAlign = align.Fill
		flexData.HGrab = true
		flexData.VGrab = true
		flexData.SizeHint.Width = 1024
		flexData.SizeHint.Height = 768
		wv.SetLayoutData(flexData)
		wv.LoadURL("https://gurpscharactersheet.com")
		content.AddChild(wv)
	} else {
		img, err := draw.AcquireImageFromFile(imagesFS, "/mountains.jpg")
		if err == nil {
			imgPanel := imagelabel.New(img)
			imgPanel.SetFocusable(true)
			_, prefSize, _ := ui.Sizes(imgPanel, layout.NoHintSize)
			imgPanel.SetSize(prefSize)
			tooltip.SetText(imgPanel, "mountains.jpg")
			scrollArea := scrollarea.New(imgPanel, scrollarea.Unmodified)
			flexData = flex.NewData()
			flexData.HAlign = align.Fill
			flexData.VAlign = align.Fill
			flexData.HGrab = true
			flexData.VGrab = true
			scrollArea.SetLayoutData(flexData)
			content.AddChild(scrollArea)

			wnd.EventHandlers().Add(event.FocusGainedType, func(evt event.Event) {
				if !installedMap[wnd] {
					crsr := getAppleCursor()
					if crsr != nil {
						imgPanel.EventHandlers().Add(event.UpdateCursorType, func(evt event.Event) {
							wnd.SetCursor(crsr)
							evt.Finish()
						})
						installedMap[wnd] = true
					}
				}
			})
		} else {
			fmt.Println(err)
		}
	}

	wnd.SetFocus(textFieldsPanel.Children()[0])
	wnd.Pack()
	wnd.ToFront()
	return wnd
}

var installedMap = make(map[*window.Window]bool)
var appleCursorOnce sync.Once
var appleCursor *cursor.Cursor

func getAppleCursor() *cursor.Cursor {
	appleCursorOnce.Do(func() {
		if img, err := draw.AcquireImageFromFile(imagesFS, "/classic-apple-logo.png"); err == nil {
			imgSize := img.Size()
			appleCursor = cursor.NewCursor(img.Data(), geom.Point{
				X: imgSize.Width / 2,
				Y: imgSize.Height / 2,
			})
		}
	})
	return appleCursor
}

func createListPanel() ui.Widget {
	list := list.New(&label.CellFactory{})
	list.Append("One",
		"Two",
		"Three with some long text to make it interesting",
		"Four",
		"Five")
	list.EventHandlers().Add(event.SelectionType, func(evt event.Event) {
		fmt.Print("Selection changed in list. Now:")
		index := -1
		first := true
		for {
			index = list.Selection.NextSet(index + 1)
			if index == -1 {
				break
			}
			if first {
				first = false
			} else {
				fmt.Print(",")
			}
			fmt.Printf(" %d", index)
		}
		fmt.Println()
	})
	list.EventHandlers().Add(event.ClickType, func(evt event.Event) {
		fmt.Println("Double-clicked on list")
	})
	_, prefSize, _ := ui.Sizes(list, layout.NoHintSize)
	list.SetSize(prefSize)
	scrollArea := scrollarea.New(list, scrollarea.Fill)
	flexData := flex.NewData()
	flexData.HAlign = align.Fill
	flexData.VAlign = align.Fill
	flexData.HGrab = true
	flexData.VGrab = true
	scrollArea.SetLayoutData(flexData)
	return scrollArea
}

func addSeparator(parent ui.Widget) {
	sep := separator.New(true)
	flexData := flex.NewData()
	flexData.HAlign = align.Fill
	sep.SetLayoutData(flexData)
	parent.AddChild(sep)
}

func createButtonsPanel() ui.Widget {
	panel := widget.NewBlock()
	lay := flow.New(panel)
	lay.HSpacing = 5
	lay.VSpacing = 5
	lay.VCenter = true

	createButton("Press Me", panel)
	createButton("Disabled", panel).SetEnabled(false)

	img, err := draw.AcquireImageFromFile(imagesFS, "/home.png")
	if err == nil {
		createImageButton(img, "Home", panel)
		createImageButton(img, "Home (disabled)", panel).SetEnabled(false)
	} else {
		fmt.Println(err)
	}

	img, err = draw.AcquireImageFromFile(imagesFS, "/classic-apple-logo.png")
	if err == nil {
		createImageButton(img, "Classic Apple Logo", panel)
		createImageButton(img, "Classic Apple Logo (disabled)", panel).SetEnabled(false)
	} else {
		fmt.Println(err)
	}

	return panel
}

func createButton(title string, panel ui.Widget) *button.Button {
	button := button.New(title)
	button.EventHandlers().Add(event.ClickType, func(evt event.Event) { fmt.Printf("The button '%s' was clicked.\n", title) })
	tooltip.SetText(button, fmt.Sprintf("This is the tooltip for the '%s' button.", title))
	panel.AddChild(button)
	return button
}

func createImageButton(img *draw.Image, name string, panel ui.Widget) *imagebutton.ImageButton {
	size := img.Size()
	size.Width /= 2
	size.Height /= 2
	button := imagebutton.NewImageButtonWithImageSize(img, size)
	button.EventHandlers().Add(event.ClickType, func(evt event.Event) { fmt.Printf("The button '%s' was clicked.\n", name) })
	tooltip.SetText(button, name)
	panel.AddChild(button)
	return button
}

func createCheckBoxPanel() ui.Widget {
	panel := widget.NewBlock()
	flex.NewLayout(panel)
	createCheckBox("Press Me", panel)
	createCheckBox("Initially Mixed", panel).SetState(checkbox.Mixed)
	createCheckBox("Disabled", panel).SetEnabled(false)
	check := createCheckBox("Disabled w/Check", panel)
	check.SetEnabled(false)
	check.SetState(checkbox.Checked)
	return panel
}

func createCheckBox(title string, panel ui.Widget) *checkbox.CheckBox {
	check := checkbox.NewCheckBox(title)
	check.EventHandlers().Add(event.ClickType, func(evt event.Event) { fmt.Printf("The checkbox '%s' was clicked.\n", title) })
	tooltip.SetText(check, fmt.Sprintf("This is the tooltip for the '%s' checkbox.", title))
	panel.AddChild(check)
	return check
}

func createRadioButtonsPanel() ui.Widget {
	panel := widget.NewBlock()
	flex.NewLayout(panel)

	group := radiobutton.NewGroup()
	first := createRadioButton("First", panel, group)
	createRadioButton("Second", panel, group)
	createRadioButton("Third (disabled)", panel, group).SetEnabled(false)
	createRadioButton("Fourth", panel, group)
	group.Select(first)

	return panel
}

func createRadioButton(title string, panel ui.Widget, group *radiobutton.Group) *radiobutton.RadioButton {
	rb := radiobutton.New(title)
	rb.EventHandlers().Add(event.ClickType, func(evt event.Event) { fmt.Printf("The radio button '%s' was clicked.\n", title) })
	tooltip.SetText(rb, fmt.Sprintf("This is the tooltip for the '%s' radio button.", title))
	panel.AddChild(rb)
	group.Add(rb)
	return rb
}

func createPopupMenusPanel() ui.Widget {
	panel := widget.NewBlock()
	flex.NewLayout(panel)

	createPopupMenu(panel, 1, "One", "Two", "Three", "", "Four", "Five", "Six")
	createPopupMenu(panel, 2, "Red", "Blue", "Green").SetEnabled(false)

	return panel
}

func createPopupMenu(panel ui.Widget, selection int, titles ...string) *popupmenu.PopupMenu {
	p := popupmenu.NewPopupMenu()
	tooltip.SetText(p, fmt.Sprintf("This is the tooltip for the PopupMenu with %d items.", len(titles)))
	for _, title := range titles {
		if title == "" {
			p.AddSeparator()
		} else {
			p.AddItem(title)
		}
	}
	p.SelectIndex(selection)
	p.EventHandlers().Add(event.SelectionType, func(evt event.Event) { fmt.Printf("The '%v' item was selected from the PopupMenu.\n", p.Selected()) })
	panel.AddChild(p)
	return p
}

func createTextFieldsPanel() ui.Widget {
	panel := widget.NewBlock()
	flex.NewLayout(panel)

	field := createTextField("First Text Field", panel)
	createTextField("Second Text Field (disabled)", panel).SetEnabled(false)
	createTextField("", panel).SetWatermark("Watermarked")
	field = createTextField("", panel)
	field.SetWatermark("Enter only numbers")
	field.EventHandlers().Add(event.ValidateType, func(evt event.Event) {
		if e, ok := evt.(*event.Validate); ok {
			for _, r := range field.Text() {
				if !unicode.IsDigit(r) {
					e.MarkInvalid()
					break
				}
			}
		}
	})

	return panel
}

func createTextField(text string, panel ui.Widget) *textfield.TextField {
	field := textfield.New()
	field.SetText(text)
	flexData := flex.NewData()
	flexData.HAlign = align.Fill
	flexData.HGrab = true
	field.SetLayoutData(flexData)
	tooltip.SetText(field, fmt.Sprintf("This is the tooltip for the '%s' text field.", text))
	panel.AddChild(field)
	return field
}

func createAboutWindow(evt event.Event) {
	if aboutWindow == nil {
		aboutWindow = window.NewWindow(geom.Point{}, window.TitledWindowMask|window.ClosableWindowMask)
		aboutWindow.EventHandlers().Add(event.ClosedType, func(evt event.Event) { aboutWindow = nil })
		aboutWindow.SetTitle("About " + cmdline.AppName)
		content := aboutWindow.Content()
		content.SetBorder(border.NewEmpty(geom.NewUniformInsets(10)))
		flex.NewLayout(content)
		title := label.NewWithFont(cmdline.AppName, font.EmphasizedSystem)
		flexData := flex.NewData()
		flexData.HAlign = align.Middle
		flexData.HGrab = true
		title.SetLayoutData(flexData)
		content.AddChild(title)
		desc := label.New("Simple app to demonstrate the\ncapabilities of the ui framework.")
		desc.SetLayoutData(flexData.Clone())
		content.AddChild(desc)
		aboutWindow.Pack()
	}
	aboutWindow.ToFront()
}

func createPreferencesWindow(evt event.Event) {
	fmt.Println("Preferences...")
}
