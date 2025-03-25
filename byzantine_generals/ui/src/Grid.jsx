import { useState, useEffect, useCallback } from 'react';
import { ReactFlow, Controls, Background, ConnectionMode, MarkerType, applyNodeChanges, applyEdgeChanges } from '@xyflow/react';

import BiDirectionalEdge from './BiDirectionalEdge';
import BiDirectionalNode from './BiDirectionalNode';

const edgeTypes = {
  bidirectional: BiDirectionalEdge,
};
 
const nodeTypes = {
  bidirectional: BiDirectionalNode,
};

const defaultNodes = [];
const defaultEdges = [];

function Grid() {
  let [nodes, setNodes] = useState(defaultNodes);
  let [edges, setEdges] = useState(defaultEdges);

  function getDecisions() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Decisions';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        if(data.decisions === undefined) {
          return;
        }

        data.decisions.forEach((decision) => {
          // add the decision to the node
          const index = nodes.findIndex((item) => item.id === decision.from);
          if (index > -1) {
            console.log("Adding decision to node", decision.from, decision.decision);

            let bgColor = '#fd4848';
            if(decision.decision === 'attack') {
              bgColor = '#36ff40';
            }

            nodes[index] = {
              ...nodes[index],
              style: {...nodes[index].style, backgroundColor: bgColor},
              data: {...nodes[index].data, decision: decision.decision}
            };
          }
        });

        // update the nodes if we have decisions
        if (data.decisions.length > 0) {
          setNodes(nodes);
        }
    });
  }

  function getNodes() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Nodes';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        let nodesToAdd = [];

        data.nodes.forEach((node, index) => {
          // does the item exist
          const found = nodes.find((item) => item.id === node.id);
          if (found) {
            return;
          }

          // set the positions of the nodes
          if (node.id === '0') {
            node.position = { x: 100, y: 100 };
          } else {
            node.position = { x: 0 + 200 * index, y: 500 };
          }

          nodesToAdd.push(node);
        });

        if (nodesToAdd.length === 0) {
          return;
        }

        const newNodes = [...nodes, ...nodesToAdd];
        setNodes(newNodes);
        nodes = newNodes;
      });
  }
  
  function getEdges() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Edges';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        if(data.edges === undefined) {
          return;
        }
        
        let edgesToAdd = [];

        // loop over the nodes and find the nodes for each edge
        // we need this to calculate the source and target
        data.edges.forEach((edge) => {
          const source = nodes.find((node) => node.id === edge.source);
          const target = nodes.find((node) => node.id === edge.target);
         
          // do we have these nodes
          if (!source || !target) {
            return;
          }
          
          // does the item exist
          const found = edges.find((item) => item.id === edge.id);
          if (found) {
            return;
          }
          
          console.log(`Adding edge from ${source.data.label} to ${target.data.label}`);

          // work out which node on the target to connect to
          if(source.position.y < target.position.y) {
            edge.targetHandle = 'top';
          }else if (source.position.x > target.position.x) {
            edge.targetHandle = 'right';
            edge.sourceHandle = 'left';
          } else {
            edge.targetHandle = 'left';
            edge.sourceHandle = 'right';
          }

          edge.markerEnd = {type: MarkerType.Arrow};
          edge.type = 'bidirectional';
          edgesToAdd.push(edge);
        });

        if (edgesToAdd.length > 0) {
          edges = [...edges, ...edgesToAdd];
          setEdges(edgesToAdd);
        }

      });
  }

  useEffect(() => {
    const nodesInterval = setInterval(getNodes, 5000);

    return () => clearInterval(nodesInterval);
  }, []);
  
  useEffect(() => {
    const edgesInterval = setInterval(getEdges, 5000);

    return () => clearInterval(edgesInterval);
  }, []);
  
  useEffect(() => {
    const decisionsInterval = setInterval(getDecisions, 5000);

    return () => clearInterval(decisionsInterval);
  }, []);

  const onNodesChange = useCallback(
    (changes) => {
      setNodes((oldNodes) => applyNodeChanges(changes, oldNodes));
    },
    [setNodes],
  );
  
  const onEdgesChange = useCallback(
    (changes) => {
      setEdges((oldEdges) => applyEdgeChanges(changes, oldEdges));
    },
    [setEdges],
  );

  return (
    <ReactFlow 
      nodes={nodes} 
      edges={edges} 
      onNodesChange={onNodesChange} 
      onEdgesChange={onEdgesChange} 
      connectionMode={ConnectionMode.Loose}
      nodeTypes={nodeTypes}
      edgeTypes={edgeTypes}
      fitView
    >
      <Background />
      <Controls />
    </ReactFlow>
  );
}
 
export default Grid;