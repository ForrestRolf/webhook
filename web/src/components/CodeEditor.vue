<script setup>
import {computed, nextTick, ref} from "vue";
import {SaveOutlined} from "@ant-design/icons-vue"

import * as monaco from 'monaco-editor'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import JsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import ShellCompletionProvider from '../monaco/shell-completion-provider.js'

monaco.languages.registerCompletionItemProvider("shell", ShellCompletionProvider(monaco))
self.MonacoEnvironment = {
    getWorker(_, label) {
        if (label === 'json') {
            return new JsonWorker()
        }
        return new EditorWorker()
    }
}
const editor = ref()
const initCodeEditor = (id, lang, initCode, readOnly) => {
    if (editor.value) return
    editor.value = monaco.editor.create(document.getElementById(id), {
        value: initCode || "",
        language: lang || "shell",
        automaticLayout: true,
        readOnly: readOnly
    });
}
const setContent = (code) => {
    monaco.editor.getEditors().at(0)?.setValue(code)
}
const setOptions = (options) => {
    monaco.editor.getEditors().at(0)?.updateOptions(options)
}
const setLanguage = (lang) => {
    monaco.editor.setModelLanguage(monaco.editor.getModels().at(0), lang)
}

const visible = ref(false)
const props = defineProps({
    id: {
        type: String,
        default() {
            return "code-editor"
        }
    },
    code: {
        type: String,
        default() {
            return ""
        }
    },
    onSave: {
        type: Function,
        default() {
            return () => {
            }
        }
    },
    readOnly: {
        type: Boolean,
        default() {
            return false
        }
    }
})

const title = computed(() => {
    const lang = editor.value?.getModel()?.getLanguageId()
    const label = {
        "shell": "Shell",
        "json": "JSON",
        "plaintext": "Text",
        "xml": "XML"
    }
    return label[lang] + (props.readOnly ? "" : " Editor")
})
const drawerConfig = computed(() => {
    if (props.readOnly) {
        return {
            width: "50%",
            height: "50%",
            placement: "right",
            closable: true,
        }
    }
    return {
        width: "100%",
        height: "100%",
        placement: "bottom",
        closable: false,
    }
})

const open = (lang, code) => {
    visible.value = true

    nextTick(() => {
        initCodeEditor(props.id, lang, code, props.readOnly)
        setLanguage(lang)
        setContent(code)
    })
}
const close = () => {
    visible.value = false
    dispose()
}
const handleSave = () => {
    visible.value = false
    props.onSave(monaco.editor.getEditors().at(0).getValue())
    dispose()
}
const dispose = () => {
    monaco.editor.getEditors().forEach(editor => editor.dispose());
    monaco.editor.getModels().forEach(model => model.dispose());
    editor.value.dispose()
    editor.value = null
}
const onClose = () => {
    close()
}

defineExpose({open, close, setContent, setOptions, setLanguage})
</script>

<template>
    <a-drawer
        :title="title"
        :placement="drawerConfig.placement"
        :closable="drawerConfig.closable"
        :open="visible"
        :width="drawerConfig.width"
        :height="drawerConfig.height"
        @close="onClose"
    >
        <div :id="props.id" style="width: 100%;height: 100%;"></div>

        <template #extra>
            <a-space>
                <a-button @click="handleSave" type="primary" v-show="!props.readOnly">
                    <template #icon>
                        <SaveOutlined/>
                    </template>
                    Save
                </a-button>
            </a-space>
        </template>
    </a-drawer>
</template>

<style scoped>

</style>
