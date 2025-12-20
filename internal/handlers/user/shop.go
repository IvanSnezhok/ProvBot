package user

import (
	"fmt"

	"provbot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ShopCategory represents a product category
type ShopCategory struct {
	ID   string
	Name map[string]string // language -> name
}

// ShopProduct represents a product
type ShopProduct struct {
	ID          string
	CategoryID  string
	Name        map[string]string // language -> name
	Description map[string]string // language -> description
	Price       float64
	OrderURL    string // External URL for ordering
}

// ShopHandler handles shop-related commands
type ShopHandler struct {
	categories []ShopCategory
	products   []ShopProduct
}

// NewShopHandler creates a new shop handler
func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		categories: getDefaultCategories(),
		products:   getDefaultProducts(),
	}
}

// getDefaultCategories returns default shop categories
func getDefaultCategories() []ShopCategory {
	return []ShopCategory{
		{
			ID: "routers",
			Name: map[string]string{
				"ua": "Роутери",
				"en": "Routers",
				"ru": "Роутеры",
			},
		},
		{
			ID: "cables",
			Name: map[string]string{
				"ua": "Кабелі та аксесуари",
				"en": "Cables & Accessories",
				"ru": "Кабели и аксессуары",
			},
		},
		{
			ID: "other",
			Name: map[string]string{
				"ua": "Інше обладнання",
				"en": "Other Equipment",
				"ru": "Другое оборудование",
			},
		},
	}
}

// getDefaultProducts returns default shop products
func getDefaultProducts() []ShopProduct {
	return []ShopProduct{
		{
			ID:         "router_tp_link",
			CategoryID: "routers",
			Name: map[string]string{
				"ua": "TP-Link Archer C6",
				"en": "TP-Link Archer C6",
				"ru": "TP-Link Archer C6",
			},
			Description: map[string]string{
				"ua": "Двохдіапазонний Wi-Fi роутер AC1200",
				"en": "Dual-band Wi-Fi router AC1200",
				"ru": "Двухдиапазонный Wi-Fi роутер AC1200",
			},
			Price:    1500,
			OrderURL: "",
		},
		{
			ID:         "router_mikrotik",
			CategoryID: "routers",
			Name: map[string]string{
				"ua": "MikroTik hAP ac lite",
				"en": "MikroTik hAP ac lite",
				"ru": "MikroTik hAP ac lite",
			},
			Description: map[string]string{
				"ua": "Компактний двохдіапазонний роутер",
				"en": "Compact dual-band router",
				"ru": "Компактный двухдиапазонный роутер",
			},
			Price:    2200,
			OrderURL: "",
		},
		{
			ID:         "patch_cord_5m",
			CategoryID: "cables",
			Name: map[string]string{
				"ua": "Патч-корд UTP 5м",
				"en": "UTP Patch cord 5m",
				"ru": "Патч-корд UTP 5м",
			},
			Description: map[string]string{
				"ua": "Мережевий кабель Cat5e, 5 метрів",
				"en": "Network cable Cat5e, 5 meters",
				"ru": "Сетевой кабель Cat5e, 5 метров",
			},
			Price:    50,
			OrderURL: "",
		},
	}
}

// HandleShowCategories shows shop categories
func (h *ShopHandler) HandleShowCategories(ctx *handlers.HandlerContext) error {
	lang := "ua"
	if ctx.User != nil {
		lang = ctx.User.Language
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range h.categories {
		name := category.Name[lang]
		if name == "" {
			name = category.Name["ua"]
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(name, "shop_cat_"+category.ID),
		))
	}

	// Add back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "shop_back"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("shop_categories"))
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleShowCategoryProducts shows products in a category
func (h *ShopHandler) HandleShowCategoryProducts(ctx *handlers.HandlerContext, categoryID string) error {
	lang := "ua"
	if ctx.User != nil {
		lang = ctx.User.Language
	}

	var products []ShopProduct
	for _, p := range h.products {
		if p.CategoryID == categoryID {
			products = append(products, p)
		}
	}

	if len(products) == 0 {
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Translator.Get("shop_no_products"))
		_, err := ctx.Bot.Send(msg)
		return err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, product := range products {
		name := product.Name[lang]
		if name == "" {
			name = product.Name["ua"]
		}
		label := fmt.Sprintf("%s - %.0f грн", name, product.Price)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, "shop_prod_"+product.ID),
		))
	}

	// Add back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "shop_categories"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Get category name
	categoryName := ""
	for _, cat := range h.categories {
		if cat.ID == categoryID {
			categoryName = cat.Name[lang]
			if categoryName == "" {
				categoryName = cat.Name["ua"]
			}
			break
		}
	}

	text := fmt.Sprintf("%s\n\n%s:", ctx.Translator.Get("shop_products_in_category"), categoryName)

	// Edit message if callback, send new if not
	if ctx.Update.CallbackQuery != nil {
		editMsg := tgbotapi.NewEditMessageText(
			ctx.Update.CallbackQuery.Message.Chat.ID,
			ctx.Update.CallbackQuery.Message.MessageID,
			text,
		)
		editMsg.ReplyMarkup = &keyboard
		_, err := ctx.Bot.Send(editMsg)
		return err
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// HandleShowProductDetails shows product details
func (h *ShopHandler) HandleShowProductDetails(ctx *handlers.HandlerContext, productID string) error {
	lang := "ua"
	if ctx.User != nil {
		lang = ctx.User.Language
	}

	var product *ShopProduct
	for _, p := range h.products {
		if p.ID == productID {
			product = &p
			break
		}
	}

	if product == nil {
		callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, ctx.Translator.Get("shop_product_not_found"))
		_, _ = ctx.Bot.Request(callbackConfig)
		return nil
	}

	name := product.Name[lang]
	if name == "" {
		name = product.Name["ua"]
	}

	description := product.Description[lang]
	if description == "" {
		description = product.Description["ua"]
	}

	text := fmt.Sprintf("<b>%s</b>\n\n%s\n\n%s: %.0f грн",
		name, description, ctx.Translator.Get("shop_price"), product.Price)

	var rows [][]tgbotapi.InlineKeyboardButton

	// Add order button (contact support to order)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("shop_order"), "shop_order_"+product.ID),
	))

	// Add back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("back"), "shop_cat_"+product.CategoryID),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	editMsg := tgbotapi.NewEditMessageText(
		ctx.Update.CallbackQuery.Message.Chat.ID,
		ctx.Update.CallbackQuery.Message.MessageID,
		text,
	)
	editMsg.ParseMode = "HTML"
	editMsg.ReplyMarkup = &keyboard
	_, err := ctx.Bot.Send(editMsg)
	return err
}

// HandleOrder handles product order
func (h *ShopHandler) HandleOrder(ctx *handlers.HandlerContext, productID string) error {
	lang := "ua"
	if ctx.User != nil {
		lang = ctx.User.Language
	}

	var product *ShopProduct
	for _, p := range h.products {
		if p.ID == productID {
			product = &p
			break
		}
	}

	if product == nil {
		return nil
	}

	name := product.Name[lang]
	if name == "" {
		name = product.Name["ua"]
	}

	// Show order confirmation message
	text := ctx.Translator.Getf("shop_order_message", name)

	callbackConfig := tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)

	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, text)
	_, err := ctx.Bot.Send(msg)
	return err
}
