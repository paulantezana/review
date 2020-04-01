import React from 'react';
import { Modal, Form, Input, Divider, InputNumber } from 'antd';

const formItemLayout = {
    labelCol: {
        xs: { span: 24 },
        sm: { span: 8 },
    },
    wrapperCol: {
        xs: { span: 24 },
        sm: { span: 16 },
    },
};

const AddForm = Form.create()(
    class extends React.Component {
        // Confirm
        handleConfirm = () => {
            const { form, onSave } = this.props;
            form.validateFields((err, values) => {
                if (err) {
                    return;
                }
                onSave(values);
                form.resetFields();
            });
        };

        // Render
        render() {
            const { form, loading, currentItem, onCancel, modalVisible } = this.props;
            const { getFieldDecorator, getFieldsValue } = form;
   
            return (
                <Modal
                    title="Programa de estudios"
                    okText="Guardar"
                    confirmLoading={loading}
                    onCancel={onCancel}
                    onOk={this.handleConfirm}
                    visible={modalVisible}
                >
                    <Form layout="horizontal" >
                        <Divider orientation="left">Instituto</Divider>

                        <Form.Item hasFeedback label="Prefijo">
                            {getFieldDecorator('prefix', {
                                initialValue: currentItem.prefix,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="Prefijo Nombre Corto">
                            {getFieldDecorator('prefix_short_name', {
                                initialValue: currentItem.prefix_short_name,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="Nombre Instituto">
                            {getFieldDecorator('institute', {
                                initialValue: currentItem.institute,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Divider orientation="left">Director(@)</Divider>

                        <Form.Item hasFeedback label="Nombre completo director(@)">
                            {getFieldDecorator('director', {
                                initialValue: currentItem.director,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="Nivel academico director(@)">
                            {getFieldDecorator('academic_level_director', {
                                initialValue: currentItem.academic_level_director,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="Nivel academico nombre corto director(@)">
                            {getFieldDecorator('short_academic_level_director', {
                                initialValue: currentItem.short_academic_level_director,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>

                        <Divider orientation="left">% Practicas modulares</Divider>

                        <Form.Item hasFeedback label="% minimo de horas de PM">
                            {getFieldDecorator('min_hours_practice_percentage', {
                                initialValue: currentItem.min_hours_practice_percentage,
                                rules: [
                                    {
                                        pattern: /^([3-9]|[1-8][0-9]|9[0-9]|10[01])$/,
                                        message: '¡Solo se permiten valores numéricos de 3 a 101!',
                                    },
                                    { required: true, message: '¡Por favor ingrese un número!' },
                                ],
                            })(<InputNumber min={3} max={255} step={1} />)}
                        </Form.Item>

                        <Divider orientation="left">Datos adicionales</Divider>

                        <Form.Item hasFeedback label="Nombre del año">
                            {getFieldDecorator('year_name', {
                                initialValue: currentItem.year_name,
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="Email">
                            {getFieldDecorator('email', {
                                initialValue: currentItem.email,
                            })(<Input />)}
                        </Form.Item>

                        <Form.Item hasFeedback label="URL Del sitio web">
                            {getFieldDecorator('web_site', {
                                initialValue: currentItem.web_site,
                            })(<Input />)}
                        </Form.Item>
                    </Form>
                </Modal>
            );
        }
    }
);

export default AddForm;