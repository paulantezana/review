(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([[17],{"2HGp":function(e,a,t){e.exports={form:"sn-pages-user-login-index-form",forgot:"sn-pages-user-login-index-forgot"}},"6T+2":function(e,a,t){"use strict";t.r(a);t("+L6B");var r=t("2/Rp"),n=(t("sRBo"),t("kaz8")),o=(t("y8nQ"),t("Vl3Y")),l=(t("5NDa"),t("5rEg")),i=(t("Pwec"),t("CtXQ")),s=t("2Taf"),m=t.n(s),c=t("vZ4D"),p=t.n(c),d=t("l4Ni"),u=t.n(d),g=t("ujKo"),f=t.n(g),h=t("rlhR"),b=t.n(h),E=t("MhPg"),v=t.n(E),w=t("q1tI"),y=t.n(w),M=t("MuoO"),k=t("mOP9"),N=t("LLXN"),O=t("2HGp"),j=t.n(O),q=function(e){function a(e){var t;return m()(this,a),t=u()(this,f()(a).call(this,e)),t.handleSubmit=t.handleSubmit.bind(b()(t)),t}return v()(a,e),p()(a,[{key:"handleSubmit",value:function(e){var a=this;e.preventDefault(),this.props.form.validateFields(function(e,t){e||a.props.dispatch({type:"login/login",payload:t})})}},{key:"render",value:function(){var e=this.props.form.getFieldDecorator,a=this.props.loading;return y.a.createElement("div",null,y.a.createElement("div",{className:j.a.form},y.a.createElement(o["a"],{onSubmit:this.handleSubmit},y.a.createElement(o["a"].Item,null,e("user_name",{rules:[{required:!0,message:Object(N["formatMessage"])({id:"validation.userName.required"})}]})(y.a.createElement(l["a"],{prefix:y.a.createElement(i["a"],{type:"user",style:{color:"rgba(0,0,0,.25)"}}),placeholder:Object(N["formatMessage"])({id:"app.login.userName"})}))),y.a.createElement(o["a"].Item,null,e("password",{rules:[{required:!0,message:Object(N["formatMessage"])({id:"validation.password.required"})}]})(y.a.createElement(l["a"].Password,{prefix:y.a.createElement(i["a"],{type:"lock",style:{color:"rgba(0,0,0,.25)"}}),placeholder:Object(N["formatMessage"])({id:"app.login.password"})}))),y.a.createElement(o["a"].Item,null,e("remember",{valuePropName:"checked",initialValue:!1})(y.a.createElement(n["a"],null,y.a.createElement(N["FormattedMessage"],{id:"app.login.remember-me"}))),y.a.createElement(k["a"],{className:j.a.forgot,to:"/user/forgot"},y.a.createElement(N["FormattedMessage"],{id:"app.login.forgot-password"})),y.a.createElement(r["a"],{type:"primary",loading:a,htmlType:"submit",block:!0},y.a.createElement(N["FormattedMessage"],{id:"app.login.login"}))))))}}]),a}(w["Component"]),x=o["a"].create()(q),F=function(e){var a=e.loading;return{loading:a.effects["login/login"]}};a["default"]=Object(M["connect"])(F)(x)}}]);