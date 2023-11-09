<script setup>
import {computed} from "vue";

const props = defineProps({
    trigger: {
        type: Object,
        default() {
            return {}
        }
    },
})
const trigger = computed({
    get() {
        return props.trigger
    }
})
</script>

<template>
    <a-space>
        <a-tag color="blue">
            {{trigger.match.parameter.source}}
            <span>{&nbsp;</span>
            <a-tag color="purple" class="no-margin-rgt">{{trigger.match.parameter.name}}</a-tag>
            <span>&nbsp;}</span>
        </a-tag>

        <template v-if="trigger.match.type === 'value'">
            <span><span v-if="trigger.match.not" class="notice">!</span>=</span>
            <span>{{trigger.match.value}}</span>
        </template>
        <template v-if="trigger.match.type === 'regex'">
            <span>~=<span v-if="trigger.match.not" class="notice">!</span></span>
            <span>{{trigger.match.regex}}</span>
        </template>
    </a-space>
</template>

<style scoped>

</style>