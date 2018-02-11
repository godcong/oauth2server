/**
* login验证
*/
function validations() {
  $.mvalidateExtend({
    username:{
      required : true,   
      pattern : /^[0-9A-Za-z.@\-\_]{6,}$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入用户名',
        pattern :'输入用户名（至少6个字符）'
      }
    },
    username2:{
      required : true,   
      pattern : /^[A-Za-z]\w{5,}$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入用户名',
        pattern :'用户名必须字母开头，至少六位'
      }
    },
    mobile:{
      required : true,   
      pattern : /^0?1[3|4|5|7|8|9][0-9]\d{8}$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入手机号',
        pattern :'手机格式错误'
      }
    },
    email:{
      required : true,   
      pattern : /^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入邮箱',
        pattern :'邮箱格式错误'
      }
    },
    password:{
      required : true,
      pattern : /^\w{6,}$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入密码',
        pattern :'密码由字母或数字或下划线组成，至少六位'
      }
    },
    code:{
      required : true,   
      pattern : /^\d{4}$/,
      each:function(){                  
      },
      descriptions:{
        required : '请输入验证码',
        pattern :'验证码为四位数字'
      }
    },
  });
}
/**
* login验证一个页面中有不同表单的切换，所以修改了验证的插件，每次验证把当前的表单id传进来
*/
function validateFormLogin(fn,form) {
  $(''+form).mvalidate({
      type:1,
      onKeyup:true,
      sendForm:true,
      firstInvalidFocus:false,
      valid:function(event,options){
        //点击提交按钮时,表单通过验证触发函数
        fn();
        event.preventDefault();
      },
      invalid:function(event, status, options){
      },
      eachField:function(event,status,options){
      },
      eachValidField:function(val){},
      eachInvalidField:function(event, status, options){},
      conditional:{
        confirmpwd2:function(val){
          var flag;
          return (val==$("#pwd2").val()) ? true :false; 
        }
      },
      descriptions:{
        username:{
          required : '请输入用户名'
        },
        password:{
          required : '请输入密码'
        },
        code:{
          required:'请输入验证码'
        },
        confirmpassword:{
          required : '请再次输入密码',
          conditional : '两次密码不一样'
        },
      },
      form:form
  });
}