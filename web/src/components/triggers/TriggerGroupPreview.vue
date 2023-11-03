<script setup>
import {computed} from "vue";
import _ from "lodash";
import TriggerPreview from "./TriggerPreview.vue";

const props = defineProps({
    triggers: {
        type: Object,
        default() {
            return {
                and: []
            }
        }
    },
    depth: {
        type: Number,
        default() {
            return 1
        }
    }
})

const triggers = computed(() => {
    return props.triggers
})
const logic = computed(() => {
    let keys = Object.keys(triggers.value)
    return keys.length > 0 ? keys[0] : "and"
})

const isTriggerGroup = (v) => {
    let keys = _.keys(v)
    return keys.includes("and") || keys.includes("or")
}

</script>

<template>
    <a-space wrap>
        <span v-if="props.depth > 1">(</span>
        <div v-for="(trigger, i) in triggers[logic]">
            <TriggerPreview v-if="!isTriggerGroup(trigger)" v-model:trigger="triggers[logic][i]"></TriggerPreview>
            <TriggerGroupPreview v-else :triggers="trigger" :depth="props.depth + 1"></TriggerGroupPreview>
            <span v-show="i < triggers[logic].length - 1" class="margin-lft">
                <a-tag>{{logic.toUpperCase()}}</a-tag>
            </span>
        </div>
        <span v-if="props.depth > 1">)</span>
    </a-space>
</template>

<style scoped>

</style>
