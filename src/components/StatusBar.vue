<template>
    <div v-if="message" class="alert" :class="alertClass" role="alert">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">×</span>
        </button>
        {{ message }}
    </div>
</template>

<script>
import { computed } from 'vue';

export default {
    props: {
        message: {
            type: String,
            default: '',
        },
        type: {
            type: String,
            default: 'info',
            validator: (value) => ['success', 'info', 'warning', 'danger'].includes(value),
        },
    },
    setup(props) {
        const alertClass = computed(() => {
            return `alert-${props.type}`;
        });

        return {
            alertClass,
        };
    },
};
</script>

<style scoped>
.alert {
    margin-bottom: 15px;
}
</style>