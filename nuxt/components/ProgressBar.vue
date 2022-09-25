<template>
  <div class="mb-4 pl-4 pr-4 pb-4 disabled">
    <vue-slider
      :value="value"
      :data="data"
      tooltip="active"
      :clickable="false"
      :marks="true"
    ></vue-slider>
  </div>
</template>

<script>
export default {
  name: "ProgressBar",
  data() {
    return {
      data: ["On sale", "Booked", "Approving", "Sold out"],
      steps: {
        booked: 2,
        paid: 3,
        fulfilled: 4,
      },
    };
  },
};
</script>
<style scoped>
.disabled {
  pointer-events: none;
}
</style>
