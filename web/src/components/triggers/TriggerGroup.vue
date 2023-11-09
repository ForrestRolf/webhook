<script setup>
import {computed} from "vue";
import _ from "lodash"
import Trigger from "./Trigger.vue";
import {PlusCircleOutlined, CloseCircleOutlined} from "@ant-design/icons-vue"
import emitter from "../../support/emitter.js";

const props = defineProps({
    triggers: {
        type: Object,
        default() {
            return {
                and: []
            }
        }
    },
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    },
    path: {
        type: Array,
        default() {
            return []
        }
    }
})
const emit = defineEmits(["update:triggers"])

const triggers = computed({
    get() {
        return props.triggers
    },
    set(t) {
        emit("update:triggers", t)
    }
})
const logic = computed({
    get() {
        let keys = Object.keys(triggers.value)
        return keys.length > 0 ? keys[0] : "and"
    },
    set(n) {
        let k = _.keys(triggers.value)[0]
        triggers.value[n] = triggers.value[k]
        delete triggers.value[k]
    }
})
const isTriggerGroup = (v) => {
    let keys = _.keys(v)
    return keys.includes("and") || keys.includes("or")
}

const addMatchRule = () => {
    triggers.value[logic.value].push({
        "match": {
            "parameter": {
                "source": "payload",
                "name": ""
            },
            "type": "value",
            "value": ""
        }
    })
}
const addGroupMathRule = () => {
    triggers.value[logic.value].push({
        "and": [
            {
                "match": {
                    "parameter": {
                        "source": "payload",
                        "name": ""
                    },
                    "type": "value",
                    "value": ""
                }
            }
        ]
    })
}
const handleRemove = (idx) => {
    triggers.value[logic.value].splice(idx, 1)
}

const generaPath = (idx) => {
    let keys = Object.keys(triggers.value[logic.value][idx])
    return _.concat(props.path, [idx, keys[0]])
}
const handleRemoveGroup = () => {
    emitter.emit("trigger-group-removed", props.path)
}
</script>

<template>
    <a-row class="trigger-group" :gutter="[12, 12]">
        <div class="remove" @click="handleRemoveGroup">
            <CloseCircleOutlined />
        </div>
        <a-col :span="24" v-for="(trigger, i) in triggers[logic]">
            <Trigger v-if="!isTriggerGroup(trigger)" v-model:trigger="triggers[logic][i]" :disabled="props.disabled"
                     @remove="handleRemove(i)"></Trigger>
            <TriggerGroup v-else :triggers="trigger" :disabled="props.disabled" :path="generaPath(i)"></TriggerGroup>

            <a-divider v-show="i < triggers[logic].length - 1">
                <a-radio-group v-model:value="logic" button-style="solid" size="small" :disabled="props.disabled">
                    <a-radio-button value="and">And</a-radio-button>
                    <a-radio-button value="or">Or</a-radio-button>
                </a-radio-group>
            </a-divider>
        </a-col>
        <a-col :span="24" class="txt-center add-more" v-show="!props.disabled">
            <a-space>
                <a-button type="text" @click="addMatchRule">
                    <template #icon>
                        <PlusCircleOutlined/>
                    </template>
                    Add match
                </a-button>
                <span>OR</span>
                <a-button type="text" @click="addGroupMathRule">
                    <template #icon>
                        <PlusCircleOutlined/>
                    </template>
                    Add group match
                </a-button>
            </a-space>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
