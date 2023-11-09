<script setup>
import {nextTick, ref} from "vue";
import {SaveOutlined} from "@ant-design/icons-vue"

import * as monaco from 'monaco-editor'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import JsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import CssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import HtmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import TsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'

self.MonacoEnvironment = {
    getWorker(_, label) {
        if (label === 'json') {
            return new JsonWorker()
        }
        if (label === 'css' || label === 'scss' || label === 'less') {
            return new CssWorker()
        }
        if (label === 'html' || label === 'handlebars' || label === 'razor') {
            return new HtmlWorker()
        }
        if (label === 'typescript' || label === 'javascript') {
            return new TsWorker()
        }
        return new EditorWorker()
    }
}
const editor = ref()
const initCodeEditor = (id, lang, initCode) => {
    if (editor.value) return
    editor.value = monaco.editor.create(document.getElementById(id), {
        value: initCode || "",
        language: lang || "shell",
        automaticLayout: true
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
    }
})
const open = (lang, code) => {
    visible.value = true

    nextTick(() => {
        initCodeEditor(props.id, lang, code)
        setLanguage(lang)
        setContent(code)
    })
}
const close = () => {
    visible.value = false
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

}

defineExpose({open, close, setContent, setOptions, setLanguage})
</script>

<template>
    <a-drawer
        title="Code Editor"
        placement="bottom"
        :closable="false"
        :open="visible"
        width="100%"
        height="100%"
        @close="onClose"
    >
        <div :id="props.id" style="width: 100%;height: 100%;"></div>

        <template #extra>
            <a-space>
                <a-button @click="handleSave" type="primary">
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
