import React from 'react';
import {
  getBezierPath,
  useStore,
  BaseEdge,
  type Edge,
  type EdgeProps,
  type ReactFlowState,
  EdgeLabelRenderer,
} from '@xyflow/react';

function EdgeLabel({ transform, label }: { transform: string; label: string }) {
  return (
    <div
      style={{
        position: 'absolute',
        background: 'transparent',
        padding: 10,
        color: '#ff5050',
        fontSize: 12,
        fontWeight: 700,
        transform,
      }}
      className="nodrag nopan"
    >
      {label}
    </div>
  );
}
 
export type GetSpecialPathParams = {
  sourceX: number;
  sourceY: number;
  targetX: number;
  targetY: number;
};
 
export const getSpecialPath = (
  { sourceX, sourceY, targetX, targetY }: GetSpecialPathParams,
  offset: number,
) => {
  const centerX = (sourceX + targetX) / 2;
  const centerY = (sourceY + targetY) / 2;
 
  return `M ${sourceX} ${sourceY} Q ${centerX} ${centerY + offset} ${targetX} ${targetY}`;
};
 
export default function CustomEdge({
  source,
  target,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  markerEnd,
  label,
}: EdgeProps<Edge<{label: string}>>) {
  const isBiDirectionEdge = useStore((s: ReactFlowState) => {
    const edgeExists = s.edges.some(
      (e) =>
        (e.source === target && e.target === source) ||
        (e.target === source && e.source === target),
    );
 
    return edgeExists;
  });
 
  const edgePathParams = {
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  };
 
  let path = '';
  let labelX = 0;
  let labelY = 0;
  let offset = '30';
 
  if (isBiDirectionEdge) {
    path = getSpecialPath(edgePathParams, sourceX < targetX ? 25 : -25);
    if (sourceX > targetX) {
      offset = '-130';
    }

    labelX = (sourceX + targetX) / 2;
    labelY = (sourceY + targetY) / 2;
  } else {
    [path, labelX, labelY] = getBezierPath(edgePathParams);
  }

  return (
    <>
      <BaseEdge path={path} markerEnd={markerEnd} />
      <EdgeLabelRenderer>
        <EdgeLabel
            transform={`translate(-50%, ${offset}%) translate(${labelX}px,${labelY}px)`}
            label={label}
          />
      </EdgeLabelRenderer>
    </>
  );
}