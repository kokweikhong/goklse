function createStockListTable() {
  new DataTable("#stock-list-table", {
    fixedColumns: true,
    fixedColumns: {
      left: 3,
    },
    paging: true,
    pageLength: 20,
    lengthMenu: [20, 40, 80, 160, 320],
    scrollCollapse: true,
    scrollX: true,
    scrollY: 300,
  });
}
