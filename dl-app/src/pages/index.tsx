import {Alert, Button, Form, Input} from 'antd';
import axios from "axios";
import {useEffect, useState} from "react";


const saveURL = async (url: string): Promise<{ err: string, task_id: string }> => {
  const res = await axios.post('http://localhost:12432/api/save', {
    url,
  })
  const {err, task_id} = res.data
  return {err, task_id}
}

const getTaskResult = async (task_id: string): Promise<{ err: string, status: 'finished' }> => {
  const res = await axios.get(`http://localhost:12432/api/get_task?task_id=${task_id}`)
  const {err, status} = res.data
  return {err, status}
}

const Demo = () => {
  const [err, setErr] = useState('')
  const [taskID, setTaskID] = useState('')
  const [loading, setLoading] = useState(false)
  const [status, setStatus] = useState('')

  useEffect(() => {
    if (!loading || !taskID) {
      return
    }

    const interval = setInterval(async () => {
      const {err, status} = await getTaskResult(taskID)
      if (err) {
        setErr(err)
        clearInterval(interval)
        return
      }
      setStatus(status)
      if (status === 'finished') {
        setLoading(false)
        clearInterval(interval)
        return
      }
    }, 1000)
  }, [taskID, loading])

  const onFinish = async (values: { url: string }) => {
    console.log('Success:', values);
    const res = await saveURL(values.url)
    setErr(res.err)
    setTaskID(res.task_id)
    setLoading(true)
  };

  return (
    <div>
      {
        err && err.length > 0 && <Alert message={err} type="error" style={{marginBottom:20}}/>
      }

      <Form
        name="basic"
        onFinish={onFinish}
        autoComplete="off"
      >
        <Form.Item
          label="链接"
          name="url"
          rules={[
            {
              required: true,
              message: '请输入链接',
            },
          ]}
        >
          <Input/>
        </Form.Item>

        <Form.Item
        >
          <Button type="primary" htmlType="submit">
            提交
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default function IndexPage() {
  return (
    <div style={{margin: 20, padding: 20}}>
      <Demo/>
    </div>
  );
}
