function circleLayout(nodeArray, w, h) {
    const centerX = w / 2;
    const centerY = h / 2;
    const radius = Math.min(w, h) / 2.5;
    const angleStep = (2 * Math.PI) / nodeArray.length;

    nodeArray.forEach((node, i) => {
        const angle = i * angleStep;
        node.x = centerX + radius * Math.cos(angle);
        node.y = centerY + radius * Math.sin(angle);
    });
}

function gridLayout(nodes, width, height) {
    const numCols = Math.ceil(Math.sqrt(nodes.length));
    const cellWidth = width / numCols;
    const cellHeight = 60; // fixed height per row

    nodes.forEach((node, i) => {
        const col = i % numCols;
        const row = Math.floor(i / numCols);

        node.x = (col + 0.5) * cellWidth;
        node.y = (row + 0.5) * cellHeight;
    });
}

function horizontalLayout(nodes, width, height) {
    const gap = width / (nodes.length + 1);
    const centerY = height / 2;

    nodes.forEach((node, i) => {
        node.x = (i + 1) * gap;
        node.y = centerY;
    });
}

function verticalLayout(nodes, width, height) {
    const gap = height / (nodes.length + 1);
    const centerX = width / 2;

    nodes.forEach((node, i) => {
        node.x = centerX;
        node.y = (i + 1) * gap;
    });
}



function drawGraphForSQLMD(canvasId, metadataMap) {
    const svg = d3.select("#" + canvasId + "-canvas");
    console.log("Drawing graph for", "#" + canvasId + "-canvas");

    const schemaObject = metadataMap[canvasId];
    console.log("Schema object", schemaObject);
    const nodes = schemaObject.dataTables.map(dt => ({
        id: dt.name,
        schema: dt.schema
    }));
    console.log("Nodes", nodes);
    const existingNames = new Set(nodes.map(n => n.id));
    console.log("Existing names", existingNames);
    function getNode(name) {
        return nodes.find(n => n.id === name);
    }
    const links = [];
    if (schemaObject.constraints) {
        schemaObject.constraints.forEach(c => {
            if (c.constraintType === "FOREIGN KEY") {
                const source = c.tableName;
                if (c.targetTable === "" || c.targetTable === null || c.targetTable === undefined) {
                    c.targetTable = eliminateSubStrings(c.constraintName, [
                        "_id_fkey", "fk_", "_id", "_fk", source + "_", source + "_id", source + "_fk", "_fkey", "_fk_id", "_" + source, source
                    ]);
                }
                const target = c.targetTable;
                if (!existingNames.has(target)) {
                    existingNames.add(target);
                    nodes.push({ id: target, schema: "unknown" });
                }
                links.push({ source, target });
            }
        });
    }

    const screenWidth = window.innerWidth;
    const screenHeight = window.innerHeight;
    const width = screenWidth * 0.85;
    const height = screenHeight * 0.85;
    console.log("Calculated width and height:", width, height);
    circleLayout(nodes, width, height);
    const container = svg.append("g").attr("class", "container");

    svg.append("defs").append("marker")
        .attr("id", "arrow")
        .attr("viewBox", "0 -5 10 10")
        .attr("refX", 12) // arrow offset from node shape
        .attr("refY", 0)
        .attr("markerWidth", 4)
        .attr("markerHeight", 4)
        .attr("orient", "auto")
        .append("path")
        .attr("d", "M0,-5L10,0L0,5")
        .attr("fill", "#999");

    container.selectAll("line.link")
        .data(links)
        .enter()
        .append("line")
        .attr("class", "graph-node-link")
        .attr("x1", d => getNode(d.source).x)
        .attr("y1", d => getNode(d.source).y)
        .attr("x2", d => getNode(d.target).x)
        .attr("y2", d => getNode(d.target).y);

    const customParam = { key: "value" }; // example parameter

    const node = container.selectAll("g.node")
        .data(nodes)
        .enter()
        .append("g")
        .attr("class", "node")
        .attr("transform", d => `translate(${d.x},${d.y})`)
        .on("mouseover", function (event, d) {
            d3.select(this).select("text")
                .transition().duration(200)
                .style("font-size", "15px");
        })
        .on("mouseout", function (event, d) {
            d3.select(this).select("text")
                .transition().duration(200)
                .style("font-size", "10px");
        })
        .call(
            d3.drag()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended)
        );
    // .on("click", handleNodeClick);

    function dragged(event, d) {
        console.log("Dragging", d.id, "to", event.x, event.y);
        d.x = event.x;
        d.y = event.y;
        d3.select(this).attr("transform", `translate(${d.x},${d.y})`);
        console.log(container)
        container.selectAll("line.link")
          .attr("x1", l => getNode(l.source).x)
          .attr("y1", l => getNode(l.source).y)
          .attr("x2", l => getNode(l.target).x)
          .attr("y2", l => getNode(l.target).y);
      }

    node.append("rect")
        .attr("class", "graph-node-rect")
        .attr("width", 150)
        .attr("height", 50)
        .attr("x", -50)
        .attr("y", -15);

    node.append("text")
        .attr("class", "graph-node-label")
        .attr("dy", 4)
        .text(d => d.id);
    const zoom = d3.zoom()
        .scaleExtent([0.5, 5]) // min & max zoom
        .on("zoom", (event) => {
            container.attr("transform", event.transform);
        });

    // Attach zoom behavior to svg
    svg.call(zoom);
}

function dragstarted(event, d) {
    d3.select(this).raise().classed("active", true);
  }
  
  function dragended(event, d) {
    d3.select(this).classed("active", false);
  }
  

function handleNodeClick(event, d) {
    console.log("Clicked on", d, d.id);
    const tableData = dvdrental.dataTables.find(t => t.name === d.id);
    if (!tableData) {
        d3.select("#table-info").html(`<p>No metadata found for ${d.id}</p>`);
        return;
    }

    // Build a small HTML list of fields:
    let fieldsHtml = "<ul>";
    tableData.fields.forEach(f => {
        fieldsHtml += `<li><b>${f.field}</b>: ${f.type}</li>`;
    });
    fieldsHtml += "</ul>";

    // Insert it into the side panel
    d3.select("#table-info").html(fieldsHtml);
}

function searchNode(){
    d3.select("#ww").on("input", function() {
        const searchTerm = this.value.toLowerCase();
        
        container.selectAll("g.node")
          .attr("opacity", d => d.id.toLowerCase().includes(searchTerm) ? 1 : 0.2);
      });
}