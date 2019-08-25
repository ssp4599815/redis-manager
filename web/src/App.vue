<template>
    <div id="app">
        <SideBar/>
        <div class="layout">
            <transition class="transition" name="slide-fade" mode="out-in" appear>
                <router-view/>
            </transition>
        </div>
    </div>
</template>
<script>
    import SideBar from '@/layout/SideBar'

    export default {
        components: {SideBar}
    }

</script>

<style lang="scss">
    @import '@/style/mixin.scss';
    @import '@/style/reset.scss';
    @import '@/style/element-custom.scss';

    $backgroud-color: #F8FAFC;

    #app {
        background: $backgroud-color;
        min-height: 100vh;
        height: 100%;
        @include slide-transition;
    }

    .layout {
        display: flex;
        flex-direction: column;
        flex: 1 1 auto;
        height: 100vh;
        overflow-x: hidden;

        & > div:last-child {
            padding: 24px;
        }
    }
</style>
