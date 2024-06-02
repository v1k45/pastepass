// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func copyTextScript() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n        function copyText(event, selector) {\n            var pastedContent = document.querySelector(selector);\n\n            // Create a range and select the text\n            var range = document.createRange();\n            range.selectNode(pastedContent);\n            window.getSelection().removeAllRanges();\n            window.getSelection().addRange(range);\n\n            // Copy the selected text\n            if (navigator.clipboard) {\n                navigator.clipboard.writeText(pastedContent.innerText).then(function() {\n                    event.dataset.tooltip = 'Copied!';\n                    event.innerText = 'Copied!';\n                    setTimeout(function() {\n                        event.innerText = 'Copy content';\n                        event.dataset.tooltip = 'Click to copy';\n                        event.blur();\n                    }, 1000);\n                });\n            }\n        }\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}