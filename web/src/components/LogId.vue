<script setup>

import {computed} from "vue";

const props = defineProps({
    id: {
        type: String
    },
    color: {
        type: String,
        default() {
            return  "#299cdb"
        }
    }
})
const style = computed(() => {
    return {
        "color": props.color
    }
})
</script>

<template>
    <span class="log-id" :style="style">[{{props.id}}]</span>
</template>

<style lang="less">
    .log-id {

    }
</style>
