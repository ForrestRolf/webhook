import useMessage from "./message.js";


export default function useClipboard() {
    const {successMessage, errorMessage} = useMessage()
    const copyToClipboard = (msg) => {
        if (navigator.clipboard && window.isSecureContext) {
            navigator.clipboard.writeText(msg).then(() => {
                successMessage("Copied").show()
            }, () => {
            })
        } else {
            const textArea = document.createElement("textarea");
            textArea.value = msg;
            textArea.style.position = "absolute";
            textArea.style.left = "-999999px";
            document.body.prepend(textArea);
            textArea.select();
            try {
                document.execCommand('copy');
            } catch (error) {
                errorMessage("Unable copy link to clipboard", error)
                alert("Unable copy link to clipboard. \n Please select manually and press ctrl+c to copy: " + msg)
            } finally {
                textArea.remove();
            }
        }
    }

    return {
        copyToClipboard
    }
}
