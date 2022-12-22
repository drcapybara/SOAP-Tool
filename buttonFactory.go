package main

//Button manufacturing facility

import (
	"github.com/gotk3/gotk3/gtk"
)

//A factory that produces the buttons and corresponding signals to be used in the view.

// adds buttons in a factory style to fixed context
func setupButtons(ctx *WindowCtx) {

	labelList := []string{"Compute Hash", "Compute Tag", "Encrypt With Password", "Decrypt With Password",
		"Generate Keypair", "Encrypt With Key", "Decrypt With Key", "Sign With Key", "Verify Signature"}

	buttonList := make([]gtk.Button, len(labelList))
	for i, label := range labelList {
		btn, _ := gtk.ButtonNewWithLabel(label)
		buttonList[i] = *btn
		ctx.fixed.Put(btn, 40, 80+i*45)
	}

	//SHA3 hash
	buttonList[0].SetTooltipMarkup("Computes a SHA3-512 hash of the text in the notepad.")
	buttonList[0].Connect("clicked", func() {
		ctx.initialState = false
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(ComputeSHA3HASH(text, ctx.fileMode))
		ctx.updateStatus("SHA3 hash computed successfully")
	})

	//HMAC
	buttonList[1].SetTooltipMarkup("Computes a keyed hash of the notepad. The resulting hash could only have been computed by parties with knowledge of the password and the message.")
	buttonList[1].Connect("clicked", func() {
		ctx.initialState = false
		password, result := passwordEntryDialog(ctx.win, "authentication")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			ctx.notePad.SetText(ComputeTaggedHash([]byte(password), []byte(text), "T"))
			ctx.updateStatus("Message tag computed successfully")
		} else {
			ctx.updateStatus("Tag computation cancelled")
		}
	})

	//Symmetric encryption
	buttonList[2].SetTooltipMarkup("Encrypts data under a passphrase. Can only be encrypted by parties with knowledge of the passphrase.")
	buttonList[2].Connect("clicked", func() {
		ctx.initialState = false
		password, result := passwordEntryDialog(ctx.win, "encryption")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			ctx.notePad.SetText(BytesToHexString(encryptPW([]byte(password), []byte(text))))
			ctx.updateStatus("Encryption successful")
		} else {
			ctx.updateStatus("Encryption cancelled")
		}
	}) //etc....

	//Symmettric Decryption
	buttonList[3].SetTooltipMarkup("Decrypts data under a passphrase. Can only be decrypted by parties with knowledge of the passphrase.")
	buttonList[3].Connect("clicked", func() {
		ctx.initialState = false
		password, result := passwordEntryDialog(ctx.win, "decryption")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			message, err := decryptPW([]byte(password), HexToBytes(text))
			if err != nil {
				ctx.updateStatus(err.Error())
			} else {
				ctx.notePad.SetText(string(message))
				ctx.updateStatus("Decryption successful")
			}
		} else {
			ctx.updateStatus("Decryption cancelled")
		}
	})

	//Keygen
	buttonList[4].SetTooltipMarkup("Generates a Schnorr E521 keypair from supplied password.")
	buttonList[4].Connect("clicked", func() {
		ctx.initialState = false
		key := KeyObj{}
		result := generateKeyPair(ctx, &key)
		if result {
			ctx.keytable.importKey(ctx, key)
		} else {
			ctx.updateStatus("Key generation cancelled")
		}
	}) //etc....

	// Keygen
	buttonList[5].SetTooltipMarkup("Encrypts data using a public key selected from the key table.")
	buttonList[5].Connect("clicked", func() {
		ctx.initialState = false
		key := KeyObj{}
		result := generateKeyPair(ctx, &key)
		if result {
			ctx.keytable.importKey(ctx, key)
		} else {
			ctx.updateStatus("Key generation cancelled")
		}
	}) //etc....

	//Reset session
	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("clicked", func() {
		showWarningDialog(ctx)
	})
	ctx.fixed.Put(reset, 40, 510)
}
