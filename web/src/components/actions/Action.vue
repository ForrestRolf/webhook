<script setup>
import {computed, ref} from "vue";
import {ArrowDownOutlined, DeleteOutlined, NodeExpandOutlined, BranchesOutlined} from "@ant-design/icons-vue"
import ShellAction from "./ShellAction.vue";
import HttpAction from "./HttpAction.vue";
import CodeEditor from "../CodeEditor.vue";
import Dispatcher from "./Dispatcher.vue";

const emit = defineEmits(["update:actions"])
const props = defineProps({
    actions: {
        type: Object,
        default() {
            return []
        }
    },
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    },
    arguments: {
        type: Array,
        default() {
            return []
        }
    }
})

const actions = computed({
    get() {
        return props.actions
    },
    set(v) {
        emit("update:actions", v)
    }
})

const components = {
    "shell": ShellAction,
    "http": HttpAction,
    "dispatcher": Dispatcher,
}
const handleRemove = (idx) => {
    actions.value.splice(idx, 1)
}
const addAction = () => {
    actions.value.push({
        "driver": "shell",
        "attributes": {
            "workingDirectory": "/tmp"
        }
    })
    emit("update:actions", actions.value)
}
const addOtherAction = (m) => {
    let attributes = {}
    switch(m.key) {
        case "http":
            attributes = {
                "method": "POST",
                "contentType": "application/json",
                "timeout": 30,
                "saveResponse": true,
            }
            break;
        case "dispatcher":
            attributes = {
                "if": {},
                "url": "",
                "method": "POST"
            }
            break;
    }
    actions.value.push({
        "driver": m.key,
        "attributes": attributes
    })
    emit("update:actions", actions.value)
}

const codeEditor = ref()
const onCodeSave = ref()
const handleCodeEditor = ({code, lang, onSave}) => {
    onCodeSave.value = onSave
    codeEditor.value.open(lang, code)
}
</script>

<template>
    <a-row :gutter="[12, 12]" class="actions">
        <a-col :span="24" v-for="(action, i) in actions">
            <a-button class="float-rgt" type="text" danger @click="handleRemove(i)" v-show="!props.disabled">
                <template #icon>
                    <DeleteOutlined/>
                </template>
            </a-button>
            <component
                :is="components[action.driver]"
                :disabled="props.disabled"
                v-model:attributes="actions[i].attributes"
                :arguments="props.arguments"
                :handle-code-editor="handleCodeEditor">
            </component>
            <a-divider v-show="i < actions.length - 1">
                <ArrowDownOutlined/>
            </a-divider>
        </a-col>
        <a-col :span="24" class="txt-center" v-show="!props.disabled">
            <a-dropdown-button @click="addAction">
                Add shell action
                <template #overlay>
                    <a-menu @click="addOtherAction">
                        <a-menu-item key="http">
                            <NodeExpandOutlined />
                            Add http action
                        </a-menu-item>
                        <a-menu-item key="dispatcher">
                            <BranchesOutlined />
                            Add dispatcher action
                        </a-menu-item>
                    </a-menu>
                </template>
            </a-dropdown-button>
        </a-col>

        <CodeEditor ref="codeEditor" :on-save="onCodeSave"></CodeEditor>
    </a-row>
</template>

<style scoped>

</style>
