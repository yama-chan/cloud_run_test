var Script = function () {


    /* initialize the external events
     -----------------------------------------------------------------*/

    $('#external-events div.external-event').each(function() {

        // create an Event Object (http://arshaw.com/fullcalendar/docs/event_data/Event_Object/)
        // it doesn't need to have a start or end
        var eventObject = {
            title: $.trim($(this).text()) // use the element's text as the event title
        };

        // store the Event Object in the DOM element so we can get to it later
        $(this).data('eventObject', eventObject);

        // make the event draggable using jQuery UI
        $(this).draggable({
            zIndex: 999,
            revert: true,      // will cause the event to go back to its
            revertDuration: 0  //  original position after the drag
        });

    });


    /* initialize the calendar
     -----------------------------------------------------------------*/

     var date = new Date();
     var d = date.getDate();
     var m = date.getMonth();
     var y = date.getFullYear();
 
     $('#calendar2').fullCalendar({
         header: {
             left: 'prev,next today',
             center: 'title',
             right: 'month,basicWeek,basicDay'
         },
         timeZone: 'local', // デフォルトは'local'  (localで良い場合は指定不要)
         editable: false,
         droppable: false, // this allows things to be dropped onto the calendar !!!
         drop: function(date, allDay) { // this function is called when something is dropped
 
             // retrieve the dropped element's stored Event Object
             var originalEventObject = $(this).data('eventObject');
 
             // we need to copy it, so that multiple events don't have a reference to the same object
             var copiedEventObject = $.extend({}, originalEventObject);
 
             // assign it the date that was reported
             copiedEventObject.start = date;
             copiedEventObject.allDay = allDay;
 
             // render the event on the calendar
             // the last `true` argument determines if the event "sticks" (http://arshaw.com/fullcalendar/docs/event_rendering/renderEvent/)
             $('#calendar').fullCalendar('renderEvent', copiedEventObject, true);
 
             // is the "remove after drop" checkbox checked?
             if ($('#drop-remove').is(':checked')) {
                 // if so, remove the element from the "Draggable Events" list
                 $(this).remove();
             }
 
         },
         events:
             {
                 //スケジュールJSONを下記のURLから取得。
                 url: '/getSchedule',
                 type: 'GET'
             },
             // [{
             // 	title: 'All Day Event!!!!!!',
             // 	start: new Date(y, m, 1)
             // },
             // {
             // 	title: 'Long Event',
             // 	start: new Date(y, m, d-5),
             // 	end: new Date(y, m, d-2)
             // },
             // {
             // 	id: 999,
             // 	title: 'Repeating Event',
             // 	start: new Date(y, m, d-3, 16, 0),
             // 	allDay: false
             // },
             // {
             // 	id: 999,
             // 	title: 'Repeating Event',
             // 	start: new Date(y, m, d+4, 16, 0),
             // 	allDay: false
             // },
             // {
             // 	title: 'Meeting',
             // 	start: new Date(y, m, d, 10, 30),
             // 	allDay: false
             // },
             // {
             // 	title: 'Lunch',
             // 	start: new Date(y, m, d, 12, 0),
             // 	end: new Date(y, m, d, 14, 0),
             // 	allDay: false
             // },
             // {
             // 	title: 'Birthday Party',
             // 	start: new Date(y, m, d+1, 19, 0),
             // 	end: new Date(y, m, d+1, 22, 30),
             // 	allDay: false
             // },
             // {
             // 	title: 'Click for Google',
             // 	start: new Date(y, m, 28),
             // 	end: new Date(y, m, 29),
             // 	url: 'http://google.com/'
             // }],
 
             //イベントの開始時間表示
             displayEventTime: true,
 
             //イベントの終了時刻を表示するようにcallbackを追加
             eventAfterRender: function(event, $el, view) {
                 var formattedTime = $.fullCalendar.formatDates(event.start, event.end, "HH:mm { - HH:mm: }");
                 // If FullCalendar has removed the title div, then add the title to the time div like FullCalendar would do
                 if($el.find(".fc-event-title").length === 0) {
                     $el.find(".fc-event-time").text(formattedTime + " - " + event.title);
                 }
                 else {
                     $el.find(".fc-event-time").text(formattedTime);
                 }
             },
             //スケジュールをクリック時の処理
             // eventClick: function(event) {
             //   // opens events in a popup window
             //   window.open(event.url, 'gcalevent', 'width=700,height=600');
             //   return false;
             // },
             eventClick: function(calEvent, jsEvent, view) {
               // ***** 今回はここにスケジュールのクリックイベントを追加 *****
               $('#modalTitle').html(calEvent.title);                  // モーダルのタイトルをセット
               // $('#modalBody').html(calEvent.description);          // モーダルの本文をセット
               $('#td-modal-title').html(calEvent.title);              // モーダルのタイトルをセット
               $('#td-modal-startdate').html(calEvent.startdatetime);  // モーダルの開始時間をセット
               $('#td-modal-enddate').html(calEvent.enddatetime);      // モーダルの終了時刻をセット
               $('#td-modal-description').html(calEvent.description);  // モーダルの説明をセット
               $('#btn-modal-edit').click(function(){                  //editボタンクリック時の遷移処理
                 if(calEvent.id == null && calEvent.id == ""){
                   window.location.href = '/editSchedule';
                 }else{
                   window.location.href = '/editSchedule?id=' + calEvent.id;
                 }  
               });
               $('#calendarModal').modal(); // モーダル着火
             },
             selectable: true,
             //スケジュールをフォーカス時のポップアップ表示処理
             eventRender: function(event, element) {
               //フォーカス時にイベントの説明を表示する。暫定処理
               element.attr('title', event.description);
             }
     });
 
}();