require ('jquery')

const Ajax = (($) => {

    class Ajax {
      constructor() {
  
      }
  
      static do(method,
                url,
                data={},
                async=true,
                success=()=>{},
                error=(jqXHR, textStatus, errorThrown)=>{
                  console.log('ajax error')
                  var err_msg = jqXHR.responseJSON['msg']
                  alert(err_msg);
                  return false;
                }) {
        method = method.toUpperCase()
        if (method != 'GET' && method != 'POST') {
          alert(`method必须为GET或POST，当前为${method}`)
          return false
        }
        console.log(`${method} - ${url}`)
        $.ajax(
          url,
          {
            dataType: 'json',
            cache: false,
            data: JSON.stringify(data),
            contentType: "application/json",
            method: method,
            processdata: false,
            async: async,
            error: error,
            success: success
          }
        )
      }
  
      static file(url,
                  data,
                  async=false,
                  success=()=>{},
                  error=(jqXHR, textStatus, errorThrown)=>{
                    var err_msg = jqXHR.responseJSON['msg']
                    alert(err_msg);
                    return false;
                  }) {
        console.log(`POST - ${url}`)
        $.ajax(
          url,
          {
            enctype: 'multipart/form-data',
            data: data,
            method: "POST",
            cache: false,
            processData: false,
            contentType: false,
            async: async,
            error: error,
            success: success
          }
        )
      }
  
      static get(url,
                 data={},
                 async=true,
                 success=()=>{},
                 error=(jqXHR, textStatus, errorThrown)=>{
                   var err_msg = jqXHR.responseJSON['msg']
                   alert(err_msg);
                   return false;
                 }) {
        Ajax.do("GET", url, data, async, success, error)
      }
  
      static post(url,
                  data={},
                  async=true,
                  success=()=>{},
                  error=(jqXHR, textStatus, errorThrown)=>{
                    console.log(jqXHR)
                    console.log(errorThrown)
                    var err_msg = jqXHR.responseJSON['msg']
                    alert(err_msg);
                    return false;
                  }) {
        Ajax.do("POST", url, data, async, success, error)
      }
  
    }
  
    return Ajax
  
  })(jQuery)
  
export default Ajax
  