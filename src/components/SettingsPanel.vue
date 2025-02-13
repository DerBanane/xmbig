<template>
  <form style="padding-bottom: 25px; margin-top: -50px">
    <label>
      <input
        id="hideOffline"
        type="checkbox"
        checked
        data-toggle="toggle"
        data-size="mini"
        data-onstyle="success"
        data-offstyle="danger"
        @change="emit('hideOfflineChanged', $event.target.checked)"
      />
      <span>Hide offline miners</span>
    </label>
    <label>
      <input
        id="showOfflineNotification"
        type="checkbox"
        checked
        data-toggle="toggle"
        data-size="mini"
        data-onstyle="success"
        data-offstyle="danger"
      />
      <span>Notify when miner went offline</span>
    </label>
    <label>
      <input
        id="groupByAlgo"
        type="checkbox"
        checked
        data-toggle="toggle"
        data-size="mini"
        data-onstyle="success"
        data-offstyle="danger"
        @change="emit('groupByAlgoChanged', $event.target.checked)"
      />
      <span>Group by Algo</span>
    </label>
    <br />
    <button id="resetClientStatusList" type="button" class="btn btn-danger fa fa-trash-o">
      Reset ClientStatusList
    </button>
  </form>
</template>

<script>
export default {
  emits: ['hideOfflineChanged', 'groupByAlgoChanged'],
  setup(props, { emit }) {
    return { emit };
  },
};
</script>