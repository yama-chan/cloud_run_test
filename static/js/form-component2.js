var Script = function() {

  //checkbox and radio btn

  var d = document;
  var safari = (navigator.userAgent.toLowerCase().indexOf('safari') != -1) ? true : false;
  var gebtn = function(parEl, child) {
    return parEl.getElementsByTagName(child);
  };
  onload = function() {


    if (!d.getElementById || !d.createTextNode) return;
    var ls = gebtn(d, 'label');
    for (var i = 0; i < ls.length; i++) {
      var l = ls[i];
      if (l.className.indexOf('label_') == -1) continue;
      var inp = gebtn(l, 'input')[0];
      if (l.className == 'label_check') {
        l.className = (safari && inp.checked == true || inp.checked) ? 'label_check c_on' : 'label_check c_off';
        l.onclick = check_it;
      };
      if (l.className == 'label_radio') {
        l.className = (safari && inp.checked == true || inp.checked) ? 'label_radio r_on' : 'label_radio r_off';
        l.onclick = turn_radio;
      };
    };
  };
  var check_it = function() {
    var inp = gebtn(this, 'input')[0];
    if (this.className == 'label_check c_off' || (!safari && inp.checked)) {
      this.className = 'label_check c_on';
      if (safari) inp.click();
    } else {
      this.className = 'label_check c_off';
      if (safari) inp.click();
    };
  };
  var turn_radio = function() {
    var inp = gebtn(this, 'input')[0];
    if (this.className == 'label_radio r_off' || inp.checked) {
      var ls = gebtn(this.parentNode, 'label');
      for (var i = 0; i < ls.length; i++) {
        var l = ls[i];
        if (l.className.indexOf('label_radio') == -1) continue;
        l.className = 'label_radio r_off';
      };
      this.className = 'label_radio r_on';
      if (safari) inp.click();
    } else {
      this.className = 'label_radio r_off';
      if (safari) inp.click();
    };
  };



  $(function() {

    // Tags Input
    $(".tagsinput").tagsInput();

    // Switch
    $("[data-toggle='switch']").wrap('<div class="switch" />').parent().bootstrapSwitch();

  });



  //color picker

  $('.cp1').colorpicker({
    format: 'hex'
  });
  $('.cp2').colorpicker();


  //date picker

  if (top.location != location) {
    top.location.href = document.location.href;
  }

  //jQuery初期表示処理　datetimepicker
  jQuery(function(){

    // datetime picker
    jQuery('#datetimepicker').datetimepicker();

    // datetime picker start
    jQuery('#date_timepicker_start').datetimepicker({
     format:'d.m.Y H:i',
     lang:'ja',
     onShow:function( ct ){
      this.setOptions({
       maxDateTime:jQuery('#date_timepicker_end').val()?jQuery('#date_timepicker_end').val():false
      })
     },
     timepicker:true
    });

    // datetime picker end
    jQuery('#date_timepicker_end').datetimepicker({
     format:'d.m.Y H:i',
     lang:'ja',
     onShow:function( ct ){
      this.setOptions({
       minDateTime:jQuery('#date_timepicker_start').val()?jQuery('#date_timepicker_start').val():false
      })
     },
     timepicker:true
    });
   });

  //javascript初期表示処理　datepicker
  $(function() {
    window.prettyPrint && prettyPrint();
    $('#dp1').datepicker({
      format: 'mm-dd-yyyy'
    });
    $('#dp2').datepicker();
    $('#dp3').datepicker();
    $('#dp3').datepicker();
    $('#dpYears').datepicker();
    $('#dpMonths').datepicker();


    var startDate = new Date(2012, 1, 20);
    var endDate = new Date(2012, 1, 25);
    $('#dp4').datepicker()
      .on('changeDate', function(ev) {
        if (ev.date.valueOf() > endDate.valueOf()) {
          $('#alert').show().find('strong').text('The start date can not be greater then the end date');
        } else {
          $('#alert').hide();
          startDate = new Date(ev.date);
          $('#startDate').text($('#dp4').data('date'));
        }
        $('#dp4').datepicker('hide');
      });
    $('#dp5').datepicker()
      .on('changeDate', function(ev) {
        if (ev.date.valueOf() < startDate.valueOf()) {
          $('#alert').show().find('strong').text('The end date can not be less then the start date');
        } else {
          $('#alert').hide();
          endDate = new Date(ev.date);
          $('#endDate').text($('#dp5').data('date'));
        }
        $('#dp5').datepicker('hide');
      });

    // disabling dates
    var nowTemp = new Date();
    var now = new Date(nowTemp.getFullYear(), nowTemp.getMonth(), nowTemp.getDate(), 0, 0, 0, 0);

    var checkin = $('#dpd1').datepicker({
      onRender: function(date) {
        return date.valueOf() < now.valueOf() ? 'disabled' : '';
      }
    }).on('changeDate', function(ev) {
      if (ev.date.valueOf() > checkout.date.valueOf()) {
        var newDate = new Date(ev.date)
        newDate.setDate(newDate.getDate() + 1);
        checkout.setValue(newDate);
      }
      checkin.hide();
      $('#dpd2')[0].focus();
    }).data('datepicker');
    var checkout = $('#dpd2').datepicker({
      onRender: function(date) {
        return date.valueOf() <= checkin.date.valueOf() ? 'disabled' : '';
      }
    }).on('changeDate', function(ev) {
      checkout.hide();
    }).data('datepicker');
  });



  //daterange picker

  //2019/10/10 yamashita 日付処理のエラー対処
  var nowTemp = new Date();
  var Last7days = new Date(nowTemp.setDate(nowTemp.getDate() - 6));
  var Last30days = new Date(nowTemp.setDate(nowTemp.getDate() - 29));
  var firstDayofmonth = new Date(nowTemp.getFullYear(), nowTemp.getMonth(), 1);
  var lastDayofmonth = new Date(nowTemp.getFullYear(), nowTemp.getMonth() + 1, 0);
  var firstDayofLastMonth = new Date(nowTemp.getFullYear(), nowTemp.getMonth() - 1, 1);
  var lastDayofLastmonth = new Date(nowTemp.getFullYear(), nowTemp.getMonth() - 1, 0);

  $('#reservation').daterangepicker();


  $('#reportrange').daterangepicker({
      ranges: {
        'Today': ['today', 'today'],
        'Yesterday': ['yesterday', 'yesterday'],
        'Last 7 Days': [Last7days, 'today'],
        'Last 30 Days': [Last30days, 'today'],
        'This Month': [firstDayofmonth, lastDayofmonth],
        'Last Month': [firstDayofLastMonth, lastDayofLastmonth]
      },
      opens: 'left',
      format: 'MM/dd/yyyy',
      separator: ' to ',
      startDate: Last30days,
      endDate: nowTemp,
      minDate: '01/01/2012',
      maxDate: '12/31/2013',
      locale: {
        applyLabel: 'Submit',
        fromLabel: 'From',
        toLabel: 'To',
        customRangeLabel: 'Custom Range',
        daysOfWeek: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
        monthNames: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
        firstDay: 1
      },
      showWeekNumbers: true,
      buttonClasses: ['btn-danger']
    },
    function(start, end) {
      $('#reportrange span').html(start.getMonth() + ' ' + start.getDate() + ', ' + start.getFullYear() + ' - ' + end.getMonth() + ' ' + end.getDate() + ', ' + end.getFullYear());
    }
  );

  //Set the initial state of the picker label
  $('#reportrange span').html(Last30days.getMonth() + ' ' + Last30days.getDate() + ', ' + Last30days.getFullYear() + ' - ' + nowTemp.getMonth() + ' ' + nowTemp.getDate() + ', ' + nowTemp.getFullYear());


}();
