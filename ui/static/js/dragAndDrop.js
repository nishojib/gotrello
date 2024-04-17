(function () {
  const addTask = (() => {
    let pending = Promise.resolve();

    const run = async (url, options) => {
      try {
        await pending;
      } catch (err) {
        alert(
          'Something went wrong while syncing. Triggering full page reload'
        );
        window.location.reload();
      } finally {
        return fetch(url, options);
      }
    };

    // update pending promise so that next task could await for it
    return (url, options) => (pending = run(url, options));
  })();

  htmx.onLoad(function (content) {
    var sortables = content.querySelectorAll('.sortable');
    for (var i = 0; i < sortables.length; i++) {
      var sortable = sortables[i];
      new Sortable(sortable, {
        animation: 150,
        group: 'shared',
        onEnd: (e) => {
          const oldColumnID = e.from.getAttribute('data-column-id');
          const newColumnID = e.to.getAttribute('data-column-id');

          const prevItemID =
            e.to.children[e.newIndex - 1]?.getAttribute('data-item-id') ?? '';
          const nextItemID =
            e.to.children[e.newIndex + 1]?.getAttribute('data-item-id') ?? '';
          const itemID =
            e.to.children[e.newIndex]?.getAttribute('data-item-id') ?? '';
          const body = new FormData();

          body.set('oldColumnID', oldColumnID);
          body.set('newColumnID', newColumnID);
          body.set('prevItemID', prevItemID);
          body.set('nextItemID', nextItemID);
          body.set('itemID', itemID);

          if (!(oldColumnID === newColumnID && e.oldIndex === e.newIndex)) {
            addTask('/move-item', {
              method: 'POST',
              body,
            });
          }
        },
      });
    }
  });
})();
