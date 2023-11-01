<script setup>

const props = defineProps({
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    }
})
</script>

<template>
    <a-row class="action-item">
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="Working directory">
                    <a-input></a-input>
                </a-form-item>
                <a-form-item label="Scripts">
                    <a-textarea :rows="8"></a-textarea>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
