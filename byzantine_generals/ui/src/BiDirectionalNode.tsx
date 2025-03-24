import React, { memo } from 'react';
import {
  type BuiltInNode,
  type NodeProps,
  Handle,
  Position,
} from '@xyflow/react';
 
 
const BiDirectionalNode = ({ data }: NodeProps<BuiltInNode>) => {
  return (
    <div>
      <Handle type="source" position={Position.Top} id="top" />
      <Handle type="source" position={Position.Left} id="left" />
      <Handle type="source" position={Position.Right} id="right" />
      {data?.label}
    </div>
  );
};
 
export default memo(BiDirectionalNode);