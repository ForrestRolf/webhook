<script setup>
import {computed, ref} from "vue";
import {ArrowDownOutlined, DeleteOutlined, NodeExpandOutlined} from "@ant-design/icons-vue"
import ShellAction from "./ShellAction.vue";
import HttpAction from "./HttpAction.vue";

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
    "http": HttpAction
}
const handleRemove = (idx) => {
    console.log(idx)
    actions.value.splice(idx, 1)
}
const addAction = () => {
    actions.value.push({
        "driver": "shell",
        "attributes": {}
    })
    emit("update:actions", actions.value)
}
const addOtherAction = (k) => {
    actions.value.push({
        "driver": "http",
        "attributes": {}
    })
    emit("update:actions", actions.value)
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
            <component :is="components[action.driver]" :disabled="props.disabled" v-model:attributes="actions[i].attributes"></component>
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
                    </a-menu>
                </template>
            </a-dropdown-button>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
