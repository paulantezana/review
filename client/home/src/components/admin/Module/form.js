import React from 'react';
import { Modal, Form, Input, Divider, Checkbox } from 'antd';

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
            const { getFieldDecorator } = form;
   
            return (
                <Modal
                    title="APP modulo"
                    okText="Guardar"
                    confirmLoading={loading}
                    onCancel={onCancel}
                    onOk={this.handleConfirm}
                    visible={modalVisible}
                >
                    <Form layout="horizontal" >
                        <Form.Item hasFeedback label="Nombre">
                            {getFieldDecorator('name', {
                                initialValue: currentItem.name,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>
                        <Form.Item hasFeedback label="Descriptión">
                            {getFieldDecorator('description', {
                                initialValue: currentItem.description,
                                rules: [{ required: true, message: '¡Campo obligatorio!' }],
                            })(<Input />)}
                        </Form.Item>
                        <Form.Item>
                            {getFieldDecorator('is_new', {
                                valuePropName: 'checked',
                                initialValue: currentItem.is_new,
                            })(<Checkbox>Es nuevo</Checkbox>)}
                        </Form.Item>
                    </Form>
                </Modal>
            );
        }
    }
);

export default AddForm;