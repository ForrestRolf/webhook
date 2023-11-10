export default function ShellCompletionProvider(monaco) {
    return {
        provideCompletionItems: (model, position) => {
            var word = model.getWordUntilPosition(position)
            var range = {
                startLineNumber: position.lineNumber,
                endLineNumber: position.lineNumber,
                startColumn: word.startColumn,
                endColumn: word.endColumn,
            }

            const keywords = [
                "if",
                "then",
                "do",
                "else",
                "elif",
                "while",
                "until",
                "for",
                "in",
                "esac",
                "fi",
                "fin",
                "fil",
                "done",
                "exit",
                "set",
                "unset",
                "export",
                "function",
                "ab",
                "awk",
                "bash",
                "beep",
                "cat",
                "cc",
                "cd",
                "chown",
                "chmod",
                "chroot",
                "clear",
                "cp",
                "curl",
                "cut",
                "diff",
                "echo",
                "find",
                "gawk",
                "gcc",
                "get",
                "git",
                "grep",
                "hg",
                "kill",
                "killall",
                "ln",
                "ls",
                "make",
                "mkdir",
                "openssl",
                "mv",
                "nc",
                "node",
                "npm",
                "ping",
                "ps",
                "restart",
                "rm",
                "rmdir",
                "sed",
                "service",
                "sh",
                "shopt",
                "shred",
                "source",
                "sort",
                "sleep",
                "ssh",
                "start",
                "stop",
                "su",
                "sudo",
                "svn",
                "tee",
                "telnet",
                "top",
                "touch",
                "vi",
                "vim",
                "wall",
                "wc",
                "wget",
                "who",
                "write",
                "yes",
                "zsh"
            ]

            var suggestions = [
                {
                    label: "ifelse",
                    kind: monaco.languages.CompletionItemKind.Snippet,
                    insertText: [
                        "if [${1:condition}]",
                        "then",
                        "\t$0",
                        "else",
                        "\t",
                        "fi",
                    ].join("\n"),
                    insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                    documentation: "If-Else Statement",
                    range: range,
                },
            ]

            for (const key of keywords) {
                if (key.startsWith(word.word)) {
                    suggestions.push({
                        label: key,
                        kind: monaco.languages.CompletionItemKind.Keyword,
                        insertText: key,
                        range: range,
                    });
                }
            }
            return {suggestions: suggestions}
        }
    }
}
